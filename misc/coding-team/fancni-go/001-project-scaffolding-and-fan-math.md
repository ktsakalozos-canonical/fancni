# Task 001: Project Scaffolding & Fan Address Logic

## Context

We are building a Kubernetes CNI plugin in Go that uses Ubuntu Fan Networking. This is the first task — setting up the project structure and implementing the core fan address computation logic that everything else depends on.

The full plan is at `misc/plan.md`.

## Objective

1. Initialize Go module and project directory structure
2. Implement fan address math (subnet, gateway, bridge name computation)
3. Implement CNI config parsing
4. Set up Makefile with build/test targets
5. Unit tests for all logic

## Scope

### Go module

- Module path: `github.com/ktsakalozos-canonical/fancni`
- Go version: 1.22+ (use latest stable)
- No external dependencies yet (this task is pure computation + stdlib)

### Directory structure to create

```
cmd/fancni/main.go          # Stub entrypoint (just a main func that exits, placeholder)
internal/fan/fan.go          # Fan address math
internal/fan/fan_test.go     # Tests
internal/config/config.go    # CNI config parsing from io.Reader
internal/config/config_test.go
Makefile
```

### internal/fan — Fan address computation

All functions are pure — no I/O, no exec, no netlink.

**`ComputeSubnet(overlayNetwork string, hostIP net.IP) (string, error)`**
- Given overlay `240.0.0.0/8` and host IP `172.16.3.4`, returns `"240.3.4.0/24"`
- Takes first octet of overlay, 3rd and 4th octets of host IP
- Validate: overlay must be valid CIDR, hostIP must be IPv4

**`ComputeGateway(overlayNetwork string, hostIP net.IP) (net.IP, error)`**
- Returns the `.1` address of the computed subnet
- e.g., `240.3.4.1` for the example above

**`ComputeBridgeName(overlayNetwork string) (string, error)`**
- Returns `"fan-<first_octet>"`, e.g., `"fan-240"` for `240.0.0.0/8`

**`ComputeUnderlayArg(hostIP net.IP, underlayPrefix int) string`**
- Returns the string used as `fanctl -u` argument, e.g., `"172.16.3.4/16"`

### internal/config — CNI config parsing

```go
type NetConfig struct {
    CNIVersion     string `json:"cniVersion"`
    Name           string `json:"name"`
    Type           string `json:"type"`
    OverlayNetwork string `json:"overlayNetwork"`
    UnderlayPrefix int    `json:"underlayPrefix"`
}
```

**`Parse(r io.Reader) (NetConfig, error)`**
- Read JSON from reader, unmarshal into NetConfig
- Validate: overlayNetwork must be valid IPv4 CIDR, underlayPrefix must be > 0
- Apply defaults: if `underlayPrefix` is 0, default to 16; if `overlayNetwork` is empty, default to `"240.0.0.0/8"`

### Makefile

```makefile
.PHONY: build test clean

build:
	go build -o _output/bin/fancni ./cmd/fancni/

test:
	go test ./... -v -count=1

clean:
	rm -rf _output/
```

### cmd/fancni/main.go

Minimal stub — will be fleshed out in Phase 4:
```go
package main

import "fmt"

func main() {
    fmt.Println("fancni CNI plugin")
}
```

## Non-goals

- No netlink, no iptables, no exec calls
- No IPAM logic
- No CNI command handling (ADD/DEL/CHECK/VERSION)
- No Dockerfiles or Helm

## Constraints

- All fan functions must be pure (no side effects) — this makes them trivially testable
- Use `net` stdlib for IP parsing, not string splitting (the PoC used string splitting which is fragile)

## Acceptance Criteria

- `go build ./...` succeeds
- `go test ./...` passes with coverage of fan math and config parsing
- Test cases must include:
  - Standard case: `240.0.0.0/8` + `172.16.3.4` → subnet `240.3.4.0/24`, gateway `240.3.4.1`, bridge `fan-240`
  - Different overlay: `10.0.0.0/8` + `172.16.5.6` → subnet `10.5.6.0/24`, gateway `10.5.6.1`, bridge `fan-10`
  - Error cases: invalid CIDR, IPv6 overlay (should error), nil hostIP
  - Config parsing: valid JSON, missing fields (defaults applied), invalid overlay CIDR
