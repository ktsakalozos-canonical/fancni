# Fancni Codebase Knowledge Reference

## Code Patterns & Idioms

### Error Handling

- Use `fmt.Errorf("context: %w", err)` for wrapping errors to enhance traceability.
- In `cmd/fancni/main.go`, errors from the main logic are handled by `writeCNIError(err)`, resulting in an exit code of 1.
- When executing external commands (e.g., `fanctl`), check for `exec.ErrNotFound` and provide actionable error messages (see `internal/fan/fanctl.go`).
- For file operations in IPAM, errors are propagated and wrapped; corrupt entries are detected and reported.

### Logging

- Logging is initialized in `cmd/fancni/main.go` to write to `/var/log/fancni.log` (append mode). If the log file can't be opened, logs are redirected to `os.Stderr`.
- The log format is plain text, utilizing Go's standard `log` package.
- Log the invocation context: `CNI_COMMAND` and `CNI_CONTAINERID`.
- Log warnings for recoverable issues (e.g., log file not available).

### Configuration

- CNI configuration is read from `stdin` using `config.Parse(os.Stdin)`.
- The configuration struct is named `NetConfig` (see `internal/config/config.go`).
- Overlay network is specified as a CIDR string (e.g., `"240.0.0.0/8"`).
- Host IP is auto-detected via a UDP dial trick (no traffic is sent).

## Naming Conventions

### Files & Packages

- Main binary: `cmd/fancni/main.go`
- Internal packages are organized as `internal/<domain>/<file>.go`:
  - `fan` — Fan networking logic
  - `ipam` — IP address management
  - `config` — CNI config parsing
  - `netutil` — netlink helpers
  - `iptables` — iptables helpers
  - `cni` — CNI plugin logic
- Test files are named `<file>_test.go` and reside in the same package directory.

### Functions & Variables

- Functions are named based on their purpose:
  - `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg` in `internal/fan/fan.go`
  - `EnsureBridge` in `internal/fan/fanctl.go`
  - `Allocate`, `Lookup`, `Free` for the IPAM interface
  - `HandleAdd`, `HandleDel`, `HandleCheck`, `HandleVersion` for CNI plugin
- Variables include:
  - `overlayNetwork`, `hostIP`, `podCIDR`, `containerID`, `dataDir`
  - IPAM state files: `ipam.json`, lock file: `ipam.lock`
- Test helpers are named `parseIP` and panic on invalid input for brevity.

## Testing Patterns & Framework Usage

- Unit tests utilize Go's standard `testing` package.
- Test files are named `<file>_test.go` and are located alongside their source files.
- Test functions are named `Test<Function>_<Scenario>`.
- Helper functions (e.g., `parseIP`) are used to reduce boilerplate code.
- Tests check for:
  - Correct output values
  - Error handling (invalid input, edge cases)
  - Idempotency (e.g., IPAM allocation)
- End-to-end tests are executed via shell scripts (`tests/e2e/test-e2e.sh`) and can be run using `make e2e`.
- E2E scripts validate multi-node scenarios, container networking, and cross-node forwarding.

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
  This command executes `go test ./... -v -count=1`.

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

## Common Gotchas

- Ensure that the log file directory exists and has the correct permissions; otherwise, logging will fail silently.
- When modifying configuration files, ensure they conform to the expected structure to avoid runtime errors.
- Be cautious with concurrent access to IPAM state files; use appropriate locking mechanisms to prevent data corruption.
- When running E2E tests, ensure that the necessary network configurations and permissions are in place to avoid failures.
- The `rock-build` command requires the Rockcraft tool to be installed and configured properly; ensure it is in your PATH.
