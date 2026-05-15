# Fancni Codebase Knowledge Reference

## Edge Cases Not Handled or Tested

### Fan Networking (internal/fan)
- **IPv6 Inputs:** The `ComputeSubnet`, `ComputeGateway`, and `validateOverlay` functions reject IPv6 overlays but only via `To4() == nil`. Improve error handling to provide explicit feedback for IPv6 inputs.
- **Overlay CIDR Validation:** Ensure that the overlay CIDR is validated to confirm it is indeed /8. Current logic only checks the first octet, which may lead to incorrect mappings.
- **Bridge Name Collisions:** The `ComputeBridgeName` function generates bridge names based solely on the first octet of the overlay. This can lead to collisions with overlays like "240.0.0.0/8" and "240.1.0.0/8". Implement a more robust naming strategy to avoid collisions.
- **fanctl exec error:** If `fanctl` is present but fails (e.g., due to invalid arguments or permissions), the error is surfaced but not retried or handled gracefully. Implement retry logic or better error handling.

### IPAM (internal/ipam/file_ipam.go)
- **Corrupt State File:** If `ipam.json` is partially written or corrupted, all IPAM operations will fail. Introduce a recovery or backup mechanism to handle this scenario.
- **ContainerID Collisions:** Validate the format of container IDs to prevent unexpected allocations from malicious or malformed IDs.
- **Exhaustion:** The `Allocate` function returns an error if all IPs in the pod subnet are allocated, but this is not surfaced or handled in higher layers. Ensure that higher layers can handle this error appropriately.
- **Concurrent Allocation:** The use of file locking (`flock`) for concurrency does not handle lock starvation or deadlock scenarios. Implement a timeout or deadlock detection mechanism.
- **Stale Allocations:** There is no TTL or garbage collection for container allocations, leading to orphaned allocations persisting indefinitely. Introduce a mechanism to clean up stale allocations.
- **IPv6 PodCIDR:** The current implementation only supports IPv4. Passing an IPv6 podCIDR will lead to silent failures or panics. Implement support for IPv6 or provide clear error messages.
- **PodCIDR Validation:** Validate that the podCIDR is a valid /24 to prevent allocation logic from breaking or assigning invalid IPs.

### CNI Plugin (internal/cni/plugin.go)
- **Missing CNI_COMMAND:** If `CNI_COMMAND` is unset, `main.go` returns "unknown CNI command" without logging the full environment for debugging. Add logging for the environment variables.
- **Malformed CNI Config:** If stdin is not valid JSON or missing required fields, the error is returned but not logged in detail. Enhance logging for better debugging.
- **HandleVersion:** The version info is returned, but compatibility with the CNI spec version in the config is not validated. Implement version compatibility checks.
- **ADD/DEL/CHECK:** There is no explicit handling for partial failures (e.g., IPAM succeeds but netlink fails). Ensure that cleanup is guaranteed in such scenarios.

### Host IP Detection (cmd/fancni/main.go)
- **Multiple Interfaces:** The UDP dial trick used to detect the host IP may pick the wrong interface if multiple interfaces exist. Provide an option for manual override or validation.
- **No Network:** If the host has no default route, detection fails. The error is surfaced but not recoverable. Implement fallback mechanisms or clearer error handling.

### Logging (cmd/fancni/main.go)
- **Log File Permissions:** If `/var/log/fancni.log` is not writable, it falls back to stderr but does not attempt to log to syslog or other fallbacks. Implement additional logging options.
- **Log Rotation:** The log file is opened in append mode without rotation or size limits, which can lead to indefinite growth. Implement log rotation or size limits.

## Race Conditions and Concurrency Hazards

### FileIPAM (internal/ipam/file_ipam.go)
- **Flock Usage:** The exclusive flock on a lock file for all operations can lead to indefinite blocking if a process crashes while holding the lock. Consider adding a timeout for lock acquisition.
- **No Lock Timeout:** There is no timeout or deadlock detection for lock acquisition. Implement these features to enhance reliability.
- **Read/Write Ordering:** While operations are serialized, lock contention may cause delays if multiple processes start simultaneously. Optimize for better performance under contention.
- **State File Write:** Writes are not atomic, which can leave `ipam.json` in a corrupt state if the process is killed during write. Implement atomic writes (e.g., write to a temp file, fsync, then rename).

### fanctl Bridge Creation (internal/fan/fanctl.go)
- **Bridge Creation Race:** Multiple instances of fancni may race to create the bridge, leading to only the first succeeding while others are no-ops. Implement robust error handling and checks to avoid this issue.
- **Exec Command:** The `fanctl` command lacks a timeout. If it hangs, fancni will also hang. Implement a timeout for command execution.

## Resource Leaks

### Goroutines
- **No goroutines are used** in the provided code, so no goroutine leaks are present.

### File Handles
- **Log File:** In `main.go`, the log file is deferred for closing, but if the process is killed, the file handle may remain open. Ensure proper cleanup on termination.
- **IPAM Lock File:** The lock file is opened for each operation and closed via defer. If the process crashes, the OS releases the lock, but the lock file remains on disk. Implement cleanup logic.

### Connections
- **UDP Dial:** The host IP detection opens a UDP connection but does not explicitly close it. While the OS should clean up, explicitly closing the connection is preferable.

## Maintainability Concerns

### Overly Complex Logic
- **Fan Mapping:** The mapping logic is simple, but the bridge name computation lacks robustness for multiple overlays. Refactor for clarity and robustness.
- **IPAM Allocation:** The allocation process is a linear scan from .2 to .254, which may not scale well for large subnets or high churn. Consider optimizing the allocation strategy.
- **Error Handling:** Errors are surfaced but not categorized. Higher layers do not distinguish between transient and permanent errors. Implement better error categorization.

### Poor Abstractions
- **IPAM Interface:** There is only one implementation (FileIPAM). While the interface is good, tests for alternate backends are missing. Add tests for different implementations.
- **CNI Plugin:** The plugin logic is monolithic, with ADD/DEL/CHECK handlers not separated for testability. Refactor into smaller, testable components.
- **Config Parsing:** The `config.Parse` function reads directly from stdin, making it less testable. Refactor to allow injection of input sources for easier testing.

### Hidden Coupling
- **PodCIDR Format:** The FileIPAM assumes that podCIDR is /24 and IPv4, which is not enforced. This leads to hidden coupling between fan and IPAM. Implement validation to enforce these assumptions.
- **Bridge Name:** The bridge name computation is tightly coupled to the overlay network format. Changes in the overlay format may break bridge logic. Decouple these components for better maintainability.

## Actionable Recommendations

### Edge Case Handling
- Implement explicit validation for overlay CIDR to ensure it is /8.
- Validate that podCIDR is IPv4 and /24 in the FileIPAM constructor.
- Introduce a backup/restore or atomic write mechanism for `ipam.json`.
- Implement TTL or garbage collection for stale container allocations.
- Improve error messages for IPv6 inputs and malformed container IDs.
- Add logic to ensure bridge name uniqueness for multiple overlays.

### Concurrency Improvements
- Add lock acquisition timeout and deadlock detection in FileIPAM.
- Make state file writes atomic (write temp, fsync, rename).
- Implement retry or timeout for `fanctl` execution commands.
- Consider a transactional approach for IPAM allocation.

### Resource Management
- Explicitly close the UDP connection after host IP detection.
- Implement log rotation or size limits for log files.
- Ensure the lock file is cleaned up after operations.

### Maintainability
- Refactor CNI plugin handlers into separate functions for ADD/DEL/CHECK.
- Make config parsing injectable for easier testing.
- Add tests for corrupt state file and lock contention scenarios.
- Document assumptions regarding bridge name computation and overlay mapping.
