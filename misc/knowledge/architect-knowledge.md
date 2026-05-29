# Fancni Architecture Reference - Updated

## System Boundaries

- **Fancni** is a CNI plugin implementing deterministic Fan Networking for Kubernetes pods.
- **Scope:** Handles pod IP allocation, bridge setup, iptables rules, and CNI lifecycle (ADD/DEL/CHECK).
- **Out-of-scope:** Does not manage node-level networking outside the fan bridge, nor orchestrate cluster-wide overlays.

## Module Responsibilities

### `cmd/fancni`
- **Entrypoint:** Main binary for CNI plugin.
- **Responsibilities:**
  - Logging to `/var/log/fancni.log`.
  - Reads CNI config from stdin.
  - Detects host IP (UDP dial trick).
  - Computes pod subnet via Fan math.
  - Instantiates file-backed IPAM.
  - Dispatches to internal CNI plugin handler for ADD/DEL/CHECK/VERSION.

### `internal/config`
- **Config Parsing:** Reads and validates CNI config JSON (from stdin).
- **Responsibilities:** Converts config into typed struct (`NetConfig`).

### `internal/fan`
- **Fan Networking Math:** Pure functions for overlay/underlay mapping.
- **Responsibilities:**
  - Compute pod subnet, gateway, bridge name, underlay arg.
  - Validate overlay CIDR.
  - Bridge creation via `fanctl` (only exec in codebase).
- **Technical Decision:** All mapping logic is deterministic and stateless.

### `internal/ipam`
- **IP Address Management:** Allocates/reclaims pod IPs.
- **Responsibilities:**
  - Interface for IPAM (`Allocate`, `Lookup`, `Free`).
  - File-backed implementation (`FileIPAM`):
    - Stores allocations in JSON file.
    - Uses exclusive file lock for concurrency.
    - Allocates sequential IPs in pod subnet (.2–.254).
    - Idempotent allocation: returns existing IP if containerID already allocated.
- **Technical Debt:** Only file-backed IPAM; no support for distributed/multi-node IPAM.

### `internal/netutil`
- **Netlink Helpers:** Utility for link existence, interface manipulation.
- **Responsibilities:** Abstracts netlink operations (used by fan bridge setup).

### `internal/iptables`
- **Iptables Management:** Sets up required rules for pod traffic.
- **Responsibilities:** Ensures rules for SNAT, forwarding, etc.
- **Technical Debt:** No abstraction for nftables; assumes iptables.

### `internal/cni`
- **CNI Plugin Logic:** Implements CNI lifecycle.
- **Responsibilities:**
  - Handles ADD/DEL/CHECK/VERSION.
  - Integrates config, IPAM, host IP, fan math, netutil, iptables.
  - Outputs CNI result JSON.
- **Technical Debt:** No error recovery for partial failures; assumes atomicity.

## Dependency Graph

- `cmd/fancni` → `internal/config`, `internal/fan`, `internal/ipam`, `internal/cni`
- `internal/cni` → `internal/config`, `internal/fan`, `internal/ipam`, `internal/netutil`, `internal/iptables`
- `internal/fan` → `internal/netutil`
- `internal/ipam` → (none)
- `internal/netutil` → `github.com/vishvananda/netlink`
- `internal/iptables` → `github.com/coreos/go-iptables`
- `internal/config` → (none)

## Architectural Decisions

- **Single exec:** Only `fanctl` is invoked via exec; all other networking is via Go libraries.
- **File-based IPAM:** Chosen for simplicity and node-local operation; not cluster-aware.
- **Stateless fan math:** Overlay/underlay mapping is pure and deterministic.
- **Logging:** All plugin operations are logged to a file for traceability.
- **Error Handling:** Errors are surfaced as CNI error JSON; plugin exits with code 1.

## Incomplete/In-progress Work

- **IPAM:** Only file-backed, no distributed or multi-node support. See TODOs for future cluster-wide IPAM.
- **Iptables:** No support for nftables; migration path not defined.
- **Observability:** Logging exists, but no metrics or tracing. See `misc/nightly-dreams/observability.md`.
- **Testing:** E2E tests are being expanded (see recent commits in `tests/e2e/test-e2e.sh` and `misc/coding-team/e2e-test`). Coverage for cross-node forwarding, containerd socket, HTTP retries, and image pull issues are actively being addressed.
- **Helm Chart:** Basic chart exists, but lacks advanced templating and validation. See `deploy/helm/fancni`.
- **Automation:** Nightly workflows for knowledge distillation and e2e testing are being developed (`.github/workflows/nightly-knowledge.yml`, `.github/workflows/nightly-dreams.yml`).
- **Modularity:** Internal packages are well-separated, but plugin logic in `internal/cni` is monolithic. See `misc/nightly-dreams/modularity.md` for refactoring plans.

## Technical Debt

- **IPAM Scalability:** File-based IPAM is not suitable for multi-node or HA scenarios.
- **Error Recovery:** No rollback or cleanup for partial failures in ADD/DEL.
- **Fanctl Dependency:** Requires `fanctl` binary in PATH; not vendored or containerized.
- **Iptables-only:** No abstraction for alternative packet filters.
- **Logging:** Log rotation and log level control not implemented.
- **Configuration:** No support for dynamic config reloads or advanced validation.
- **Testing:** E2E tests are shell-based and brittle; need migration to Go-based integration tests.
- **Documentation:** User-facing docs are minimal; see `README.md` and `ARCHITECTURE.md`.

## Actionable Recommendations

1. **Expand IPAM:** Implement a distributed IPAM solution to support multi-node environments.
2. **Error Handling Improvements:** Introduce error recovery mechanisms for CNI operations to handle partial failures gracefully.
3. **Migrate to nftables:** Plan and implement a migration path to support nftables alongside iptables.
4. **Enhance Observability:** Integrate metrics and tracing capabilities to monitor plugin performance and health.
5. **Refactor CNI Logic:** Break down the monolithic `internal/cni` package into smaller, more manageable components.
6. **Improve Helm Chart:** Add advanced templating and validation to the Helm chart for better deployment flexibility.
7. **Automate Testing:** Transition E2E tests from shell scripts to Go-based integration tests for better reliability and maintainability.
8. **Documentation Update:** Enhance user-facing documentation to provide clearer guidance on installation, configuration, and troubleshooting.
