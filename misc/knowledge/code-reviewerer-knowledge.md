# Fancni Codebase Knowledge Reference

## Edge Cases Not Handled or Tested

### Fan Networking (internal/fan)
- **IPv6 Inputs:** The `ComputeSubnet`, `ComputeGateway`, and `validateOverlay` functions reject IPv6 overlays and host IPs, but only via `To4() == nil`. There is no explicit error for IPv6 input; error messages could be clearer.
- **Overlay CIDR Validation:** Only the first octet is used for overlay; if the overlay network is not /8, the mapping logic may be incorrect or misleading. No validation that overlay CIDR is actually /8.
- **Bridge Name Collisions:** `ComputeBridgeName` uses only the overlay first octet, so overlays like "240.0.0.0/8" and "240.1.0.0/8" would produce the same bridge name ("fan-240"). This could cause collisions if multiple overlays are used.
- **fanctl exec error:** If `fanctl` is present but fails (e.g., due to invalid args or permissions), the error is surfaced but not retried or handled gracefully.

### IPAM (internal/ipam/file_ipam.go)
- **Corrupt State File:** If `ipam.json` is partially written or corrupted, all IPAM operations will fail. There is no recovery or backup mechanism.
- **ContainerID Collisions:** No validation on containerID format; malicious or malformed IDs could cause unexpected allocations.
- **Exhaustion:** If all IPs in the pod subnet are allocated, `Allocate` returns an error, but this is not surfaced or handled in higher layers.
- **Concurrent Allocation:** Uses file locking (`flock`) for concurrency, but does not handle lock starvation or deadlock if lock file is left open by a crashed process.
- **Stale Allocations:** No TTL or garbage collection for container allocations; orphaned allocations persist indefinitely.
- **IPv6 PodCIDR:** Only IPv4 is supported; passing IPv6 podCIDR will cause silent failures or panics.
- **PodCIDR Validation:** No check that podCIDR is a valid /24; if not, allocation logic may break or assign invalid IPs.

### CNI Plugin (internal/cni/plugin.go)
- **Missing CNI_COMMAND:** If `CNI_COMMAND` is unset, main.go returns "unknown CNI command" but does not log the full environment for debugging.
- **Malformed CNI Config:** If stdin is not valid JSON or missing required fields, error is returned but not logged in detail.
- **HandleVersion:** Returns version info, but does not validate compatibility with CNI spec version in config.
- **ADD/DEL/CHECK:** No explicit handling for partial failures (e.g., IPAM succeeds but netlink fails); cleanup is not guaranteed.

### Host IP Detection (cmd/fancni/main.go)
- **Multiple Interfaces:** Uses UDP dial trick to detect host IP, but if multiple interfaces exist, may pick the wrong one. No override or validation.
- **No Network:** If host has no default route, detection fails; error is surfaced but not recoverable.

### Logging (cmd/fancni/main.go)
- **Log File Permissions:** If `/var/log/fancni.log` is not writable, falls back to stderr, but does not attempt to log to syslog or other fallback.
- **Log Rotation:** Log file is opened in append mode, but no rotation or size limit; log file can grow indefinitely.

## Race Conditions and Concurrency Hazards

### FileIPAM (internal/ipam/file_ipam.go)
- **Flock Usage:** Uses exclusive flock on a lock file for all operations. If a process crashes while holding the lock, other processes may block indefinitely.
- **No Lock Timeout:** There is no timeout or deadlock detection for lock acquisition.
- **Read/Write Ordering:** All operations are serialized, but if multiple processes are started simultaneously, lock contention may cause delays.
- **State File Write:** Writes are not atomic; if the process is killed during write, `ipam.json` may be left in a corrupt state.

### fanctl Bridge Creation (internal/fan/fanctl.go)
- **Bridge Creation Race:** If multiple fancni instances run simultaneously, they may race to create the bridge. Only the first succeeds; others are no-ops, but error handling is not robust.
- **Exec Command:** No timeout on `fanctl` command; if it hangs, fancni will hang.

## Resource Leaks

### Goroutines
- **No goroutines are used** in the provided code, so no goroutine leaks are present.

### File Handles
- **Log File:** main.go defers logFile.Close(), but if process is killed, file handle may remain open.
- **IPAM Lock File:** Lock file is opened for each operation and closed via defer; if process crashes, OS releases lock, but lock file remains on disk.

### Connections
- **UDP Dial:** Host IP detection opens a UDP connection, but does not explicitly close it. OS should clean up, but explicit close is preferable.

## Maintainability Concerns

### Overly Complex Logic
- **Fan Mapping:** The mapping logic is simple and deterministic, but the bridge name computation is not robust for multiple overlays.
- **IPAM Allocation:** Allocation is linear scan from .2 to .254; not scalable for large subnets or high churn.
- **Error Handling:** Errors are surfaced but not categorized; higher layers do not distinguish between transient and permanent errors.

### Poor Abstractions
- **IPAM Interface:** Only one implementation (FileIPAM); interface is good, but no tests for alternate backends.
- **CNI Plugin:** Plugin logic is monolithic; ADD/DEL/CHECK handlers are not separated for testability.
- **Config Parsing:** config.Parse reads from stdin directly; not easily testable or injectable.

### Hidden Coupling
- **PodCIDR Format:** FileIPAM assumes podCIDR is /24 and IPv4; this is not enforced or validated, leading to hidden coupling between fan and ipam.
- **Bridge Name:** Bridge name computation is tightly coupled to overlay network format; changing overlay format may break bridge logic.

## Actionable Recommendations

### Edge Case Handling
- Add explicit validation for overlay CIDR to ensure /8 network.
- Validate podCIDR is IPv4 and /24 in FileIPAM constructor.
- Add backup/restore or atomic write for `ipam.json` (e.g., write to temp file then rename).
- Add TTL or garbage collection for stale container allocations.
- Improve error messages for IPv6 inputs and malformed containerIDs.
- Add bridge name uniqueness logic for multiple overlays.

### Concurrency Improvements
- Add lock acquisition timeout and deadlock detection in FileIPAM.
- Make state file writes atomic (write temp, fsync, rename).
- Add retry or timeout for `fanctl` exec command.
- Consider using a transactional approach for IPAM allocation.

### Resource Management
- Explicitly close UDP connection after host IP detection.
- Consider log rotation or log file size limit.
- Ensure lock file is cleaned up after operations.

### Maintainability
- Refactor CNI plugin handlers into separate functions for ADD/DEL/CHECK.
- Make config parsing injectable for easier testing.
- Add tests for corrupt state file and lock contention scenarios.
- Document bridge name computation and overlay mapping assumptions.

### Testing
- Add tests for:
  - Exhausted IPAM (all IPs allocated)
  - Corrupt `ipam.json`
  - IPv6 podCIDR and overlay
  - Multiple overlays with same first octet
  - Race conditions in bridge creation and IPAM allocation
  - Host IP detection with multiple interfaces

---

## Summary Table

| Area                | Concern                                   | Recommendation                      |
|---------------------|-------------------------------------------|-------------------------------------|
| Fan Networking      | Overlay CIDR not validated as /8          | Validate overlay CIDR               |
| Fan Networking      | Bridge name collisions                    | Use overlay CIDR in bridge name     |
| IPAM                | State file corruption                     | Atomic write, backup/restore        |
| IPAM                | Lock starvation/deadlock                  | Lock timeout, deadlock detection    |
| IPAM                | No garbage collection                     | Add TTL or cleanup for allocations  |
| IPAM                | PodCIDR not validated                     | Validate podCIDR as /24 IPv4        |
| CNI Plugin          | Monolithic handler logic                  | Refactor handlers                   |
| CNI Plugin          | No partial failure cleanup                | Add cleanup logic                   |
| Logging             | No log rotation                           | Add rotation or size limit          |
| Host IP Detection   | Multiple interfaces not handled           | Allow override, improve detection   |
| Testing             | Edge cases/races not covered              | Add relevant tests                  |

---

## Critical Patterns to Watch For

- **File-based state management:** Always use atomic writes and lock timeouts.
- **Bridge naming:** Avoid collisions by including more overlay info.
- **Error propagation:** Distinguish between transient and permanent errors; log details.
- **Resource cleanup:** Always close files/connections explicitly.
- **Testing:** Cover exhaustion, corruption, race, and edge input scenarios.

---

## Reference: Key Functions and Files

- `internal/fan/fan.go`: Fan address mapping logic.
- `internal/fan/fanctl.go`: Bridge creation via exec.
- `internal/ipam/file_ipam.go`: File-backed IPAM, lock and state management.
- `cmd/fancni/main.go`: Entrypoint, logging, config parsing, host IP detection.
- `internal/cni/plugin.go`: CNI plugin command dispatch.
- `internal/config/config.go`: CNI config parsing.

---

## Checklist for Review

- [ ] Overlay CIDR validated as /8
- [ ] PodCIDR validated as /24 IPv4
- [ ] Bridge name uniqueness for multiple overlays
- [ ] Atomic file writes for IPAM state
- [ ] Lock timeout and deadlock detection for IPAM
- [ ] Explicit resource cleanup (files, connections)
- [ ] Log file rotation or size limit
- [ ] Host IP detection robust to multiple interfaces
- [ ] Tests for edge cases, races, and corruption
- [ ] Refactored handler logic for maintainability

---

**Apply this checklist and recommendations to all new code, tests, and refactors.**
