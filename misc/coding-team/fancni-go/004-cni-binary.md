# Task 004: CNI Binary — ADD/DEL/CHECK/VERSION

## Context

Phases 1-3 built the building blocks: fan address math, netlink helpers, iptables helpers, IPAM. This task wires them together into the actual CNI binary that kubelet invokes.

CNI spec version: **1.0.0**. Reference: https://github.com/containernetworking/cni/blob/spec-v1.0.0/SPEC.md

## Objective

Implement the `internal/cni/plugin.go` orchestration layer and the `cmd/fancni/main.go` entrypoint so the binary works as a CNI plugin.

## Scope

### internal/cni/plugin.go

The plugin struct holds the runtime context and orchestrates ADD/DEL/CHECK/VERSION.

```go
type Plugin struct {
    config   config.NetConfig
    ipam     ipam.IPAM
    hostIP   net.IP
    
    // From CNI env vars
    command     string
    ifName      string
    netNS       string
    containerID string
}
```

**Constructor:**
```go
func NewPlugin(cfg config.NetConfig, ipam ipam.IPAM, hostIP net.IP) *Plugin
```
Reads CNI env vars: `CNI_COMMAND`, `CNI_IFNAME`, `CNI_NETNS`, `CNI_CONTAINERID`.

**HandleAdd() error**
1. Compute fan subnet, gateway, bridge name from config + hostIP (using `internal/fan`)
2. Ensure fan bridge exists: check with `netutil.LinkExists(bridgeName)`. If missing, exec `fanctl up -o <overlay> -u <hostIP>/<prefix>` (see fanctl wrapper below)
3. Allocate IP via IPAM
4. Generate veth names: host side `vethXXXX` (random 4-char hex suffix), container side is `CNI_IFNAME`
5. Create veth pair via `netutil.CreateVethPair`
6. Attach host veth to fan bridge via `netutil.AttachToBridge`
7. Move container veth to netns via `netutil.MoveToNetNS` (open ns fd from `CNI_NETNS`)
8. Configure interface inside netns via `netutil.ConfigureInterface` (IP + default route via gateway)
9. Get MAC via `netutil.GetLinkMAC`
10. Build and print CNI Result JSON to stdout

**HandleDel() error**
1. Lookup container IP in IPAM
2. If not found, return nil (idempotent)
3. Free IP in IPAM
4. Delete host veth via `netutil.DeleteLink` (pod side auto-deleted)
   - Need to find the host veth name. Two options: (a) store it in IPAM, or (b) derive it. Since veth names are random, we can't derive them. Instead: scan links whose master is the fan bridge and whose peer is in the container's netns. Simpler approach: just skip explicit link deletion — when the container netns is destroyed by the runtime, both veth ends are automatically cleaned up. This is what most CNIs do (Flannel, Calico).
   - So: just free the IPAM entry. The veth cleanup happens automatically.

**HandleCheck() error**
1. Lookup container IP in IPAM — if not found, return error
2. Compute expected fan bridge name
3. Verify interface config via `netutil.VerifyInterfaceConfig(netNS, ifName, expectedIP/24)`
4. Return nil if everything checks out, error otherwise

**HandleVersion() error**
1. Print JSON to stdout:
```json
{"cniVersion":"1.0.0","supportedVersions":["1.0.0"]}
```

### CNI Result JSON (for ADD)

CNI 1.0.0 result format:
```json
{
  "cniVersion": "1.0.0",
  "interfaces": [
    {
      "name": "eth0",
      "mac": "aa:bb:cc:dd:ee:ff",
      "sandbox": "/var/run/netns/xxx"
    }
  ],
  "ips": [
    {
      "address": "240.3.4.2/24",
      "gateway": "240.3.4.1",
      "interface": 0
    }
  ]
}
```

Note: CNI 1.0.0 removed the `version` field from the `ips` array (it was in 0.3.1).

### internal/fan/fanctl.go — fanctl exec wrapper

Add to the existing `internal/fan` package:

**`EnsureBridge(bridgeName, overlayNetwork string, hostIP net.IP, underlayPrefix int) error`**
1. Check if bridge exists via `netutil.LinkExists(bridgeName)`
2. If exists, return nil
3. If not, exec: `fanctl up -o <overlayNetwork> -u <hostIP>/<underlayPrefix>`
4. Return error if fanctl fails

This is the ONLY exec call in the entire codebase.

### cmd/fancni/main.go — Entrypoint

Replace the stub with the real entrypoint:

1. Set up logging to `/var/log/fancni.log` (append mode)
2. Read CNI config from stdin via `config.Parse(os.Stdin)`
3. Detect host IP via `net.Dial("udp", "8.8.8.8:80")` (same trick as PoC — gets the IP of the default-route interface without actually sending traffic)
4. Compute pod CIDR via `fan.ComputeSubnet`
5. Create IPAM: `ipam.NewFileIPAM("/var/lib/cni/fancni", podCIDR)`
6. Create Plugin: `cni.NewPlugin(cfg, ipam, hostIP)`
7. Switch on `CNI_COMMAND` env var → call appropriate handler
8. On error: write CNI error JSON to stdout and exit 1

CNI error format:
```json
{"cniVersion":"1.0.0","code":100,"msg":"error message"}
```

### Tests — internal/cni/plugin_test.go

Testing the full plugin requires root + netns, so most tests will need `t.Skip`. However:
- Test `HandleVersion` output (no root needed)
- Test that the result JSON structure is valid (mock/unit test the JSON marshaling)

## Non-goals

- No Dockerfiles, no Helm, no init container
- No iptables setup in the CNI binary (that's the init container's job)
- No network policy support

## Constraints

- The only `exec.Command` call is to `fanctl` in `internal/fan/fanctl.go`
- Error output must follow CNI spec: JSON to stdout with `code` and `msg` fields
- Must handle the case where `fanctl` is not in PATH (clear error message)
- Veth host-side naming: use random suffix to avoid collisions. `fmt.Sprintf("veth%04x", rand)` with `crypto/rand` for the hex bytes
