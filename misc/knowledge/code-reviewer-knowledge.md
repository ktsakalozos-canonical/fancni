## Knowledge Reference for `fancni` Code Review

This document distills actionable knowledge for reviewing the `fancni` codebase, focusing on test coverage gaps, security-sensitive paths, code smells, and error handling issues. It is organized by component and file.

---

## 1. Test Coverage Gaps

### 1.1. `internal/fan/`
- **fan.go**:
  - **Test Coverage**: 
    - **Tested**: `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg`, error paths for invalid overlay/host IP.
    - **Untested**: Edge cases in `validateOverlay` (e.g., malformed CIDR, non-IPv4 overlays) are not fully exercised.
      - **Action**: Ensure all error returns are exercised in tests.

- **fanctl.go**:
  - **NOT unit tested**: `EnsureBridge` is not covered by tests (due to `exec.Command` call).
    - **Action**: Add tests with `exec.Command` mocking (e.g., using `testexec` or similar) to verify error handling for:
      - Bridge already exists (noop path).
      - `fanctl` not found.
      - `fanctl` returns error output.
      - Success path.

### 1.2. `internal/ipam/`
- **file_ipam.go**:
  - **Partial coverage**: Many code paths (e.g., file locking, JSON corruption, allocation exhaustion, idempotency, error propagation from file I/O) are not fully covered.
    - **Action**: Add tests for:
      - Corrupt or missing state files.
      - Lock file cannot be acquired.
      - No free IPs in subnet.
      - Double allocation (idempotency).
      - ContainerID with invalid/corrupt mapping.
      - Directory creation errors.
  - **Security**: Test that the file permissions for state and lock files are restrictive (not world-writable).

### 1.3. `internal/cni/`
- **plugin.go**:
  - **Test Coverage**:
    - **Tested**: Basic plugin logic.
    - **Untested**: Error propagation from underlying IPAM, bridge, and netlink operations.
      - **Action**: Add tests for:
        - IPAM errors (allocation, lookup, free).
        - Bridge creation failures.
        - Netlink failures (e.g., link not found, route add fails).
        - Invalid/missing CNI environment variables.

### 1.4. `cmd/fancni/main.go`
- **NOT tested**: The CLI entrypoint is not covered by tests.
  - **Action**: Add integration tests (or refactor for testability) to cover:
    - Log file open failures.
    - All CNI_COMMAND variants (ADD, DEL, CHECK, VERSION, unknown).
    - Host IP detection failures.
    - CNI config parse failures.
    - Plugin error propagation.

---

## 2. Security-Sensitive Code Paths

### 2.1. `internal/fan/fanctl.go`
- **`exec.Command` call**: Only place in codebase that shells out.
  - **Risks**:
    - Path injection: Overlay network and host IP are passed as arguments. If these are attacker-controlled, could result in command injection.
    - **Mitigation**: Validate/sanitize all arguments passed to `fanctl`.
    - **Action**: Add explicit validation for `overlayNetwork`, `hostIP`, and `underlayPrefix` before passing to `exec.Command`.
    - **Test**: Add tests for malicious/invalid input.

### 2.2. `internal/ipam/file_ipam.go`
- **File I/O and locking**:
  - **Risks**:
    - Race conditions if lock is not properly acquired.
    - File permission escalation if files are world-writable.
    - Corrupt state file could cause denial of service.
  - **Action**:
    - Ensure all file writes use restrictive permissions (0600).
    - On file corruption, fail gracefully and do not panic.
    - Validate all data read from disk.
    - Test for lock acquisition failures (e.g., locked by another process).

### 2.3. `cmd/fancni/main.go`
- **Environment variables**:
  - **Risks**:
    - Trusts CNI_COMMAND, CNI_CONTAINERID, etc. from environment.
    - **Mitigation**: Validate these variables before use.
  - **Action**:
    - Add validation for CNI_COMMAND (must be one of ADD, DEL, CHECK, VERSION).
    - Validate CNI_CONTAINERID is a reasonable identifier (alphanumeric, not empty).

---

## 3. Code Smells & Inconsistencies

### 3.1. Logging
- **cmd/fancni/main.go**:
  - Logs to `/var/log/fancni.log` by default, falls back to stderr.
    - **Issue**: If running as non-root, may not have permission to write to `/var/log`.
    - **Action**: Consider making log path configurable or fallback to `/tmp/fancni.log` if `/var/log` is not writable.

### 3.2. Error Handling
- Many error messages do not include enough context (e.g., which file, which operation).
  - **Action**: Enhance error messages to include context for easier debugging and tracing.
