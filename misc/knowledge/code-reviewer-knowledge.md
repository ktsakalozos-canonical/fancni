# Knowledge Reference for fancni Code Review

This document distills actionable knowledge for reviewing the `fancni` codebase, focusing on test coverage gaps, security-sensitive paths, code smells, and error handling issues. It is organized by component and file.

---

## 1. Test Coverage Gaps

### 1.1. `internal/fan/`
- `fan.go`:
  - **Tested**: `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg`, error paths for invalid overlay/host IP.
  - **Untested**: Some edge cases in `validateOverlay` (e.g., malformed CIDR, non-IPv4 overlays) are covered, but ensure all error returns are exercised.
- `fanctl.go`:
  - **NOT unit tested**: `EnsureBridge` is not covered by tests (due to `exec.Command` call).
    - **Action**: Add tests with `exec.Command` mocking (e.g., using `testexec` or similar) to verify error handling for:
      - Bridge already exists (noop path).
      - `fanctl` not found.
      - `fanctl` returns error output.
      - Success path.

### 1.2. `internal/ipam/`
- `file_ipam.go`:
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
- `plugin.go`:
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
- `cmd/fancni/main.go`:
  - Logs to `/var/log/fancni.log` by default, falls back to stderr.
    - **Issue**: If running as non-root, may not have permission to write to `/var/log`.
    - **Action**: Consider making log path configurable or fallback to `/tmp/fancni.log` if `/var/log` is not writable.

### 3.2. Error Handling
- Many error messages do not include enough context (e.g., which file, which operation).
  - **Action**: Always wrap errors with context using `%w` and descriptive messages.
- Some error returns are not checked (e.g., in file I/O, lock acquisition).
  - **Action**: Audit all error returns, especially in `file_ipam.go` and `plugin.go`.

### 3.3. File Permissions
- `internal/ipam/file_ipam.go`:
  - Uses `os.MkdirAll` and writes state/lock files.
    - **Action**: Ensure all files are created with `0600` permissions, not default.

### 3.4. Hardcoded Paths
- `/var/lib/cni/fancni` and `/var/log/fancni.log` are hardcoded.
  - **Action**: Make configurable via environment variable or config file.

---

## 4. Error Handling Gaps

### 4.1. `internal/fan/fanctl.go`
- If `exec.Command` fails, error output is returned, but not all error types are differentiated (e.g., permission denied, invalid arguments).
  - **Action**: Improve error messages to distinguish between not found, permission denied, and command failure.

### 4.2. `internal/ipam/file_ipam.go`
- If state file is corrupt, returns error, but does not attempt recovery.
  - **Action**: Consider backup/restore or at least log a clear message for operators.
- If no IP is available, error message should include the subnet and containerID.

### 4.3. `cmd/fancni/main.go`
- If host IP detection fails, error message is generic.
  - **Action**: Include more context (e.g., which network interface, what error).

---

## 5. Specific Patterns to Flag

### 5.1. `exec.Command` Usage
- Only in `internal/fan/fanctl.go`.
- **Flag**: Any new usage of `exec.Command` elsewhere should be scrutinized for input validation and error handling.

### 5.2. File/Lock Handling
- All stateful operations in `internal/ipam/file_ipam.go` must:
  - Check for errors on open, read, write, close.
  - Use `defer` for unlocks and file closes.
  - Not panic on corrupt or missing files.

### 5.3. Environment Variable Usage
- In `cmd/fancni/main.go`, all CNI_* env vars should be validated before use.
- **Flag**: Any new code that trusts environment variables without validation.

---

## 6. Recommendations for Review

- **Test new error paths**: Whenever new error handling is added, ensure it is covered by tests (unit or integration).
- **Security review for all file and command execution code**.
- **Validate all external input**: Overlay network, host IP, containerID, environment variables.
- **Consistent error wrapping**: Use `%w` and always provide context.
- **Restrictive file permissions**: All state and lock files must be `0600`.
- **No panics on user input or file errors**: Always return errors, never panic.

---

## 7. Summary Table

| Area                        | Weak Coverage | Security Risk | Error Handling | Action Needed                        |
|-----------------------------|--------------|---------------|---------------|--------------------------------------|
| fanctl.go (exec.Command)    | Yes          | Yes           | Partial       | Add tests, validate/sanitize input   |
| file_ipam.go (file/lock)    | Yes          | Yes           | Partial       | Add tests, check permissions         |
| plugin.go (CNI logic)       | Partial      | Low           | Partial       | Add error propagation tests          |
| main.go (entrypoint)        | Yes          | Low           | Partial       | Add integration tests, validate env  |
| Logging/Paths               | N/A          | Low           | N/A           | Make configurable, check permissions |

---

## 8. Example Test Cases to Add

- `fanctl.go`: Simulate `fanctl` not found, permission denied, invalid args, already exists.
- `file_ipam.go`: Corrupt state file, lock file busy, allocation exhaustion, idempotent allocation.
- `plugin.go`: IPAM returns error, bridge creation fails, missing env vars.
- `main.go`: All CNI_COMMAND variants, log file unwritable, host IP detection fails.

---

## 9. Example Security Checks

- Validate all arguments to `exec.Command` are not attacker-controlled or contain shell metacharacters.
- Ensure all state files are `0600` and not symlinks.
- Never trust environment variables for file paths or commands without validation.

---

**Reviewers: Use this checklist and reference to guide code reviews and to identify areas needing additional tests, validation, or security hardening.**
