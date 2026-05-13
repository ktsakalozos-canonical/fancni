# Task 002: Netlink Helpers

## Context

Phase 1 established fan address math and config parsing. This task builds the netlink abstraction layer that the CNI binary will use in Phase 4 to set up pod networking. We also add `github.com/coreos/go-iptables` for iptables management.

The full plan is at `misc/plan.md`. The project structure is at the repo root.

## Objective

Create `internal/netutil/` with functions that use `vishvananda/netlink` and `vishvananda/netns` to perform all networking operations the CNI needs. Also create `internal/iptables/` for iptables rule management via `go-iptables`.

## Scope

### Dependencies to add

```
github.com/vishvananda/netlink
github.com/vishvananda/netns
github.com/coreos/go-iptables
```

Run `go get` for each.

### internal/netutil/netutil.go — Netlink operations

All functions operate via the netlink library. No shelling out to `ip`.

**`CreateVethPair(hostVethName, containerVethName string) error`**
- Creates a veth pair using `netlink.LinkAdd` with `netlink.Veth`

**`AttachToBridge(ifName, bridgeName string) error`**
- Gets the bridge link by name, gets the interface link by name
- Sets the interface master to the bridge via `netlink.LinkSetMaster`
- Sets the interface up via `netlink.LinkSetUp`

**`MoveToNetNS(ifName string, nsFd int) error`**
- Gets the link by name
- Calls `netlink.LinkSetNsFd(link, nsFd)`

**`ConfigureInterface(nsPath, ifName, ipCIDR, gatewayIP string) error`**
- Opens the target network namespace using `netns.GetFromPath(nsPath)`
- Switches to that namespace (save current ns first, defer restore)
- In the target ns:
  - Get the link by name
  - Set link up
  - Parse ipCIDR and add address via `netlink.AddrAdd`
  - Parse gateway and add default route via `netlink.RouteAdd` with `Dst: 0.0.0.0/0, Gw: gateway`
- Restore original namespace

**`DeleteLink(ifName string) error`**
- Gets link by name, calls `netlink.LinkDel`
- If link not found, return nil (idempotent)

**`GetLinkMAC(nsPath, ifName string) (string, error)`**
- Opens target namespace
- Gets link by name
- Returns `link.Attrs().HardwareAddr.String()`

**`LinkExists(ifName string) bool`**
- Calls `netlink.LinkByName`, returns true if no error

**`VerifyInterfaceConfig(nsPath, ifName, expectedIPCIDR string) error`**
- Opens target namespace
- Verifies link exists and is up
- Verifies expected IP is assigned
- Returns error describing what's wrong, or nil if OK
- This is used by the CHECK command

### internal/iptables/iptables.go — Iptables rule management

Uses `github.com/coreos/go-iptables/iptables`.

**`EnsureForwardRules(podCIDR string) error`**
- Ensures these two rules exist in the FORWARD chain of the filter table:
  - `-s <podCIDR> -j ACCEPT`
  - `-d <podCIDR> -j ACCEPT`
- Use `AppendUnique` to avoid duplicates

**`EnsureMasqueradeRule(podSubnet, bridgeName string) error`**
- Ensures this rule exists in POSTROUTING of the nat table:
  - `-s <podSubnet> ! -o <bridgeName> -j MASQUERADE`
- Use `AppendUnique`

### Minor fix from review

In `internal/config/config.go`, add validation: reject `underlayPrefix` values > 32 or < 1 (after default application). This is a one-liner bounds check.

### Unit tests

- `internal/netutil/netutil_test.go`: Test what's feasible without root/netns. At minimum, test `LinkExists` returning false for a nonexistent link. For functions requiring root + netns, add test stubs with `t.Skip("requires root")` — they'll be validated during integration.
- `internal/iptables/iptables_test.go`: Similar — tests requiring iptables access should `t.Skip("requires root")`.
- `internal/config/config_test.go`: Add test case for `underlayPrefix: 33` → expect error.

## Non-goals

- No CNI command handling
- No IPAM
- No fanctl exec calls
- No Dockerfiles

## Constraints

- Zero calls to `exec.Command("ip", ...)` — everything through netlink
- Namespace switching must always restore the original ns (defer pattern)
- All functions must handle the "already exists" / "not found" cases gracefully (idempotent where possible)
