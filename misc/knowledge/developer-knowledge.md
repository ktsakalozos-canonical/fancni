# Fancni Codebase Knowledge Reference

## Code Patterns & Idioms

### Error Handling

- Always wrap errors with context using `fmt.Errorf("context: %w", err)` for traceability.
- In `cmd/fancni/main.go`, errors from the main logic are handled by `writeCNIError(err)` and cause exit code 1.
- When calling external commands (e.g., `fanctl`), check for `exec.ErrNotFound` and provide actionable error messages (see `internal/fan/fanctl.go`).
- For file operations (IPAM), errors are propagated and wrapped; corrupt entries are detected and reported.

### Logging

- Logging is initialized in `main.go` to `/var/log/fancni.log` (append mode). If the log file can't be opened, logs fall back to `os.Stderr`.
- Log format is plain text, using Go's standard `log` package.
- Log the invocation context: `CNI_COMMAND` and `CNI_CONTAINERID`.
- Log warnings for recoverable issues (e.g., log file not available).

### Configuration

- CNI config is read from `stdin` using `config.Parse(os.Stdin)`.
- Config struct is named `NetConfig` (see `internal/config/config.go`).
- Overlay network is specified as a CIDR string (e.g., `"240.0.0.0/8"`).
- Host IP is auto-detected via UDP dial trick (no traffic sent).

## Naming Conventions

### Files & Packages

- Main binary: `cmd/fancni/main.go`
- Internal packages: `internal/<domain>/<file>.go`
  - `fan` — Fan networking logic
  - `ipam` — IP address management
  - `config` — CNI config parsing
  - `netutil` — netlink helpers
  - `iptables` — iptables helpers
  - `cni` — CNI plugin logic

- Test files: `<file>_test.go` in the same package directory.

### Functions & Variables

- Functions are named for their purpose:
  - `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg` in `fan.go`
  - `EnsureBridge` in `fanctl.go`
  - `Allocate`, `Lookup`, `Free` for IPAM interface
  - `HandleAdd`, `HandleDel`, `HandleCheck`, `HandleVersion` for CNI plugin
- Variables:
  - `overlayNetwork`, `hostIP`, `podCIDR`, `containerID`, `dataDir`
  - IPAM state files: `ipam.json`, lock file: `ipam.lock`
- Test helpers are named `parseIP` and panic on invalid input for brevity.

## Testing Patterns & Framework Usage

- Unit tests use Go's standard `testing` package.
- Test files are named `<file>_test.go` and reside alongside their source.
- Test functions are named `Test<Function>_<Scenario>`.
- Use helper functions (e.g., `parseIP`) to reduce boilerplate.
- Tests check for:
  - Correct output values
  - Error handling (invalid input, edge cases)
  - Idempotency (e.g., IPAM allocation)
- End-to-end tests are shell scripts (`tests/e2e/test-e2e.sh`) and are invoked via `make e2e`.
- E2E scripts test multi-node scenarios, container networking, and cross-node forwarding.

## Build & Run Commands

### Build

- Build the main binary:
  ```sh
  make build
  ```
  Output: `_output/bin/fancni`

### Test

- Run all unit tests:
  ```sh
  make test
  ```
  Runs `go test ./... -v -count=1`

### Clean

- Remove build artifacts:
  ```sh
  make clean
  ```

### Helm Commands

- Generate Helm templates:
  ```sh
  make helm-template
  ```
- Lint Helm charts:
  ```sh
  make helm-lint
  ```

### End-to-End Tests

- Run end-to-end tests:
  ```sh
  make e2e
  ```

## Common Gotchas

- Ensure that the log file path (`/var/log/fancni.log`) is writable; otherwise, logs will be printed to `os.Stderr`.
- When modifying CNI configurations, ensure that the expected CIDR format is followed to avoid parsing errors.
- Be cautious with the `exec` package; always check for `exec.ErrNotFound` to handle missing executables gracefully.
- The IPAM state files (`ipam.json`, `ipam.lock`) should not be manually edited, as this can lead to inconsistencies.
