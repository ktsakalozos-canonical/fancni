# FanCNI — Go-based Kubernetes CNI using Fan Networking

## Overview

A Kubernetes CNI plugin built entirely in Go that leverages Ubuntu Fan Networking for pod-to-pod communication across nodes. Fan networking uses a mathematical mapping between an underlay network (the node IPs, typically a `/16`) and an overlay network (the pod IPs, a `/8`) to deterministically assign each node a `/24` pod subnet and route cross-node traffic via VXLAN tunnels — with no distributed state or per-peer route management required.

The system has three deliverables:
1. **CNI binary** — a Go binary placed at `/opt/cni/bin/fancni`
2. **Init container** — prepares each node (installs `ubuntu-fan`, creates fan bridge, sets iptables rules)
3. **Helm chart** — deploys everything as a DaemonSet

Go module: `github.com/ktsakalozos-canonical/fancni`

---

## Architecture

### How Fan Networking Works (recap)

Given:
- Underlay network: `172.16.0.0/16` (the node IP space)
- Overlay network: `240.0.0.0/8` (the pod IP space)
- Node IP: `172.16.3.4`

Fan computes: node's pod subnet = `240.3.4.0/24` (overlay first octet + underlay 3rd and 4th octets).

Cross-node routing is automatic: a packet destined for `240.5.6.25` is encapsulated by the fan VXLAN device and sent to underlay IP `172.16.5.6`. No FDB entries, no per-peer routes, no node watcher needed.

### Components and Interactions

```
+-----------------------------------------------------------------------------------+
|  Helm Chart (DaemonSet)                                                           |
|                                                                                   |
|  +---------------------------+    +-------------------------------------------+   |
|  | Init Container            |    | Install Container                         |   |
|  |                           |    |                                           |   |
|  | 1. chroot /host           |    | 1. Copy fancni binary to                 |   |
|  |    apt-get install        |    |    /opt/cni/bin/fancni (host mount)       |   |
|  |    ubuntu-fan             |    |                                           |   |
|  | 2. Load kernel modules    |    | 2. Write CNI config to                   |   |
|  |    (vxlan, ip_tunnel)     |    |    /etc/cni/net.d/10-fancni.conflist      |   |
|  | 3. fanctl up -o overlay   |    |    (from ConfigMap)                       |   |
|  |    -u underlay            |    |                                           |   |
|  | 4. iptables FORWARD rules |    +-------------------------------------------+   |
|  | 5. iptables MASQUERADE    |                                                    |
|  +---------------------------+    +-------------------------------------------+   |
|                                   | Pause Container (keeps pod alive)         |   |
|                                   +-------------------------------------------+   |
+-----------------------------------------------------------------------------------+

             |
             | On each node, kubelet invokes the CNI binary per pod:
             v

+-----------------------------------------------------------------------------------+
|  CNI Binary (/opt/cni/bin/fancni)                                                 |
|                                                                                   |
|  ADD:                                                                             |
|    1. Read CNI config from stdin                                                  |
|    2. Compute fan subnet from overlay + host IP                                   |
|    3. Ensure fan bridge exists (call fanctl if not)                                |
|    4. Allocate IP from file-based IPAM (flock for concurrency)                    |
|    5. Create veth pair via netlink                                                |
|    6. Attach host-side veth to fan bridge via netlink                              |
|    7. Move pod-side veth to container netns via netlink                            |
|    8. Assign IP + default route inside container netns via netlink                 |
|    9. Return CNI result JSON to stdout                                            |
|                                                                                   |
|  DEL:                                                                             |
|    1. Look up container IP in IPAM                                                |
|    2. Free IP                                                                     |
|    3. Delete host-side veth via netlink (pod-side auto-deleted)                    |
|                                                                                   |
|  CHECK:                                                                           |
|    1. Verify container veth exists and has expected IP                             |
|    2. Verify host-side veth is attached to fan bridge                              |
|    3. Return error if anything is inconsistent                                    |
|                                                                                   |
|  VERSION:                                                                         |
|    1. Return supported CNI versions JSON                                          |
+-----------------------------------------------------------------------------------+
```

### Key Design Decisions

1. **No node watcher / no daemon**: Fan VXLAN computes cross-node routes mathematically. No distributed state needed.

2. **`fanctl` via exec**: We call `fanctl` to create the fan bridge + VXLAN tunnel. This is the only exec call. All other networking operations use `vishvananda/netlink` and `coreos/go-iptables`.

3. **File-based IPAM with flock**: Each node stores a `containerID -> IP` mapping in a JSON file under `/var/lib/cni/fancni/`. Concurrent access is guarded by `syscall.Flock`. This is the same pattern used by the standard `host-local` IPAM plugin.

4. **Init container for host setup**: The init container mounts the host rootfs at `/host` and uses `chroot` to install `ubuntu-fan` via `apt-get`. This is standard practice (Calico, Cilium do the same). Kernel modules are loaded via `modprobe` in the host namespace.

5. **Two init containers**: Separation of concerns — one for host OS setup (fan package, iptables), one for CNI binary/config installation. The second is a simple `cp` from the image.

---

## Project Structure

```
fancni/
  cmd/
    fancni/                  # CNI binary entrypoint
      main.go
  internal/
    cni/                     # CNI command handlers (ADD/DEL/VERSION)
      plugin.go
      plugin_test.go
    ipam/                    # File-based IPAM implementation
      file_ipam.go
      file_ipam_test.go
    fan/                     # Fan address math + fanctl wrapper
      fan.go
      fan_test.go
    netutil/                 # Netlink helpers (veth, bridge, addr, route)
      netlink.go
      netlink_test.go
    config/                  # CNI config parsing
      config.go
      config_test.go
  deploy/
    helm/
      fancni/
        Chart.yaml
        values.yaml
        templates/
          daemonset.yaml
          configmap.yaml
          serviceaccount.yaml
          clusterrole.yaml
          clusterrolebinding.yaml
    docker/
      Dockerfile.cni         # Builds the CNI binary image (copies binary)
      Dockerfile.init        # Init container image (Ubuntu-based, has apt)
    scripts/
      init-node.sh           # Node init script (called by init container)
  go.mod
  go.sum
  Makefile
  misc/
    plan.md                  # This file
```

### Package Responsibilities

| Package | Purpose | Key Dependencies |
|---------|---------|-----------------|
| `cmd/fancni` | Entrypoint: parse CNI env vars, dispatch to handler | — |
| `internal/cni` | ADD/DEL/CHECK/VERSION logic, orchestrates IPAM + netlink calls | `internal/ipam`, `internal/netutil`, `internal/fan` |
| `internal/ipam` | Allocate/free/lookup IPs from a JSON file with flock | — |
| `internal/fan` | Fan address math (subnet, gateway, bridge name) + fanctl exec | — |
| `internal/netutil` | Thin wrappers around `vishvananda/netlink` for veth, bridge, addr, route operations in namespaces | `vishvananda/netlink`, `vishvananda/netns` |
| `internal/config` | Parse CNI JSON config from stdin | — |

### External Go Dependencies

| Dependency | Purpose |
|-----------|---------|
| `github.com/vishvananda/netlink` | Veth creation, bridge operations, IP assignment, route management |
| `github.com/vishvananda/netns` | Network namespace switching for container-side operations |
| `github.com/coreos/go-iptables/iptables` | Iptables rule management (FORWARD, MASQUERADE) |

---

## CNI Configuration Format

```json
{
  "cniVersion": "1.0.0",
  "name": "fancni",
  "type": "fancni",
  "overlayNetwork": "240.0.0.0/8",
  "underlayPrefix": 16
}
```

- `overlayNetwork`: The `/8` overlay network (default: `240.0.0.0/8`)
- `underlayPrefix`: The underlay prefix length (default: `16`). The host IP is auto-detected; this just specifies the mask.

The per-node pod subnet is computed at runtime: `overlay_first_octet.host_3rd_octet.host_4th_octet.0/24`.

---

## Init Container: Node Setup

The init container runs as a privileged container with the following host mounts:
- `/` (host root) mounted at `/host` (for chroot package installation)
- `/opt/cni/bin` for CNI binary installation
- `/etc/cni/net.d` for CNI config

### init-node.sh logic

```
1. Check if fanctl exists on host (via chroot /host which fanctl)
   - If not: chroot /host apt-get update && apt-get install -y ubuntu-fan
2. Load kernel modules: modprobe vxlan, modprobe ipip (via nsenter into host pid ns)
3. Detect host IP (default route interface)
4. Run fanctl up -o <overlay> -u <host_ip>/<underlay_prefix> (via chroot)
5. Set iptables rules:
   - FORWARD ACCEPT for pod CIDR (source)
   - FORWARD ACCEPT for pod CIDR (destination)
   - MASQUERADE for pod subnet going to non-fan-bridge interfaces
6. Copy CNI binary from /fancni to /host/opt/cni/bin/fancni
7. Write CNI config to /host/etc/cni/net.d/10-fancni.conf (from env/configmap)
```

---

## IPAM Details

**Storage**: `/var/lib/cni/fancni/ipam.json`
**Lock file**: `/var/lib/cni/fancni/ipam.lock`

Format:
```json
{
  "containerID1": "240.3.4.2",
  "containerID2": "240.3.4.3"
}
```

Allocation: iterate from `.2` to `.254` in the node's `/24`, skip `.0` (network) and `.1` (gateway/bridge). First unallocated IP wins.

Concurrency: `syscall.Flock` with `LOCK_EX` on the lock file. Retry with backoff if contended.

---

## Helm Chart Values

```yaml
overlay:
  network: "240.0.0.0/8"      # Overlay network CIDR

underlay:
  prefix: 16                   # Underlay prefix length

fan:
  encapsulation: "vxlan"       # vxlan or ipip
  dhcp: false                  # Enable dnsmasq on fan bridge

images:
  cni:
    repository: ghcr.io/ktsakalozos-canonical/fancni
    tag: latest
  init:
    repository: ghcr.io/ktsakalozos-canonical/fancni-init
    tag: latest

cniConfig:
  name: "fancni"
  cniVersion: "1.0.0"
```

---

## Implementation Phases

### Phase 1: Project scaffolding & fan address logic
- Go module init (`github.com/ktsakalozos-canonical/fancni`)
- Directory structure, Makefile with `build`, `test`, `lint` targets
- `internal/fan`: fan address computation (subnet, gateway, bridge name)
- `internal/config`: CNI config parsing from stdin
- Unit tests for fan math and config parsing

### Phase 2: Netlink helpers
- `internal/netutil`: create veth pair, attach to bridge, move to netns, assign IP, add route — all via netlink
- Unit tests (where feasible without root/netns)

### Phase 3: IPAM
- `internal/ipam`: file-based IPAM with flock
- Allocate, lookup, free operations
- Unit tests

### Phase 4: CNI binary (ADD/DEL/VERSION)
- `internal/cni`: plugin struct, ADD/DEL/CHECK/VERSION handlers that orchestrate fan + IPAM + netlink
- `cmd/fancni`: entrypoint (env var dispatch, logging, fanctl binary check)
- `fanctl` exec wrapper in `internal/fan`
- CHECK command: verify veth exists, IP assigned, host-side attached to bridge
- Integration: build binary, manually test in a k8s environment

### Phase 5: Init container & Dockerfiles
- `deploy/scripts/init-node.sh`
- `deploy/docker/Dockerfile.cni` (multi-stage: build Go binary, copy to minimal image)
- `deploy/docker/Dockerfile.init` (Ubuntu-based, includes apt, modprobe)
- Test: docker build both images

### Phase 6: Helm chart
- DaemonSet with init containers + pause container
- ConfigMap for CNI config
- ServiceAccount + RBAC (minimal, only if K8s API access is needed)
- `values.yaml` with overlay/underlay/image config
- Test: `helm template` validation, `helm install` on test cluster

### Phase 7: End-to-end validation & docs
- Test pod-to-pod on same node
- Test pod-to-pod across nodes
- Test pod-to-internet
- Basic README with quickstart
- Troubleshooting notes (fan-cleanup, common issues)

---

## Constraints & Assumptions

- **Ubuntu nodes only**: `ubuntu-fan` package is Ubuntu-specific
- **IPv4 only**: Fan networking is IPv4 only
- **Overlay `/8`, underlay `/16`**: Current `fanctl` limitation (per manpage: "Currently Fan can only apply overlay addresses with a /8 network mask, and underlay addresses with a /16 network mask")
- **Privileged init container**: Required for host package installation, kernel module loading, and iptables modification
- **No network policies**: Out of scope
- **No dual-stack**: Out of scope
- **`fanctl` is the only exec call**: All other operations use Go libraries (netlink, go-iptables)

## Out of Scope

- IPv6 / dual-stack support
- Network policy enforcement
- Non-Ubuntu distro support
- Reimplementing fanctl / VXLAN logic
- Dynamic node watcher (fan VXLAN handles cross-node routing mathematically)
- Windows node support
- Custom IPAM backends (DHCP, etcd, etc.)
