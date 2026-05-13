# Fancni Architecture Reference

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

- **IPAM:** Plan migration to distributed IPAM for multi-node support.
- **Observability:** Add Prometheus metrics and structured logging.
- **Error Handling:** Implement rollback/cleanup for partial failures.
- **Testing:** Expand Go-based integration tests; reduce reliance on shell scripts.
- **Helm Chart:** Add validation, templating, and upgrade hooks.
- **Fanctl:** Consider vendoring or containerizing fanctl for portability.
- **Iptables:** Abstract iptables logic for future nftables support.
- **Config:** Add schema validation and support for dynamic reloads.

## References to In-progress/Planning

- `misc/nightly-dreams/*`: Contains architectural dreams and planning for scalability, modularity, automation, observability, maintainability, integration, validation, performance, security, testing, prioritization.
- `misc/coding-team/e2e-test/*`: Tracks e2e test improvements, bug fixes, and coverage expansion.
- `.github/workflows/nightly-knowledge.yml`: Knowledge distillation workflow for automated architectural documentation.
- `.github/workflows/nightly-dreams.yml`: Nightly architectural planning and prioritization.

---

## Internal Package Summary

| Package         | Responsibility                                   | Key Files                     |
|-----------------|--------------------------------------------------|-------------------------------|
| `cmd/fancni`    | Entrypoint, logging, CNI dispatch                | `main.go`                     |
| `internal/config` | CNI config parsing/validation                  | `config.go`, `config_test.go` |
| `internal/fan`  | Fan math, bridge setup, overlay/underlay mapping | `fan.go`, `fanctl.go`         |
| `internal/ipam` | IP allocation, file-backed state                 | `ipam.go`, `file_ipam.go`     |
| `internal/netutil` | Netlink helpers                              | `netutil.go`                  |
| `internal/iptables` | Iptables rules management                    | `iptables.go`                 |
| `internal/cni`  | CNI plugin logic (ADD/DEL/CHECK/VERSION)         | `plugin.go`                   |

---

## Key Architectural Patterns

- **Pure Functions:** Fan math is stateless and testable.
- **File Locking:** IPAM uses exclusive file locks for concurrency.
- **Idempotency:** IPAM allocation is idempotent per containerID.
- **Error Propagation:** All errors bubble up to CNI error JSON.
- **Single Responsibility:** Internal packages are focused and decoupled.

---

## Areas for Immediate Attention

- **Distributed IPAM:** Needed for multi-node Kubernetes clusters.
- **Observability:** Metrics and structured logs for production readiness.
- **Testing:** Robust integration tests for all CNI lifecycle events.
- **Fanctl Portability:** Ensure fanctl is available in all deployment environments.
- **Helm Chart:** Improve for production deployment (validation, upgrades).

---

## Recent Architectural Changes

- **E2E Testing:** Major expansion in shell-based e2e tests; coverage for cross-node, containerd, HTTP retries, image pulls.
- **Knowledge Distillation:** Automated workflows for architectural documentation.
- **Modularity Planning:** Nightly dreams and planning files outline future refactoring and modularity improvements.

---

## TODOs & Planning References

- See `misc/nightly-dreams/*` for modularity, scalability, observability, automation, integration, maintainability, validation, performance, security, testing, prioritization.
- See `misc/coding-team/e2e-test/*` for e2e test improvements and bugfixes.
- See `.github/workflows/nightly-knowledge.yml` for automated architectural knowledge updates.

---
