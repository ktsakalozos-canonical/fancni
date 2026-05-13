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

### Naming Conventions

#### Files & Packages

- Main binary: `cmd/fancni/main.go`
- Internal packages: `internal/<domain>/<file>.go`
  - `fan` — Fan networking logic
  - `ipam` — IP address management
  - `config` — CNI config parsing
  - `netutil` — netlink helpers
  - `iptables` — iptables helpers
  - `cni` — CNI plugin logic

- Test files: `<file>_test.go` in the same package directory.

#### Functions & Variables

- Functions are named for their purpose:
  - `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg` in `fan.go`
  - `EnsureBridge` in `fanctl.go`
  - `Allocate`, `Lookup`, `Free` for IPAM interface
  - `HandleAdd`, `HandleDel`, `HandleCheck`, `HandleVersion` for CNI plugin
- Variables:
  - `overlayNetwork`, `hostIP`, `podCIDR`, `containerID`, `dataDir`
  - IPAM state files: `ipam.json`, lock file: `ipam.lock`
- Test helpers are named `parseIP` and panic on invalid input for brevity.

### Testing Patterns & Framework Usage

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

### Build & Run Commands

#### Build

- Build the main binary:
  ```sh
  make build
  ```
  Output: `_output/bin/fancni`

#### Test

- Run all unit tests:
  ```sh
  make test
  ```
  Runs `go test ./... -v -count=1`

#### Clean

- Remove build artifacts:
  ```sh
  make clean
  ```

#### Docker

- Build CNI plugin image:
  ```sh
  make docker-build-cni
  ```
- Build init container image:
  ```sh
  make docker-build-init
  ```
- Build both:
  ```sh
  make docker-build
  ```

#### Helm

- Render Helm templates:
  ```sh
  make helm-template
  ```
- Lint Helm chart:
  ```sh
  make helm-lint
  ```

#### E2E

- Run end-to-end tests:
  ```sh
  make e2e
  ```
  (Invokes `bash tests/e2e/test-e2e.sh`)

### Common Gotchas

#### Fan Networking

- Overlay network must be IPv4 CIDR (e.g., `"240.0.0.0/8"`). IPv6 overlays are rejected.
- Host IP must be IPv4. IPv6 addresses are not supported.
- Fan subnet calculation is deterministic: overlay first octet + underlay 3rd/4th octets.
- Gateway IP is always `.1` of the subnet.

#### IPAM

- File-backed IPAM uses exclusive file locks (`ipam.lock`) for concurrency safety.
- Allocated IPs are stored in `ipam.json` as a map of containerID → IP string.
- Allocations are idempotent: repeated calls for the same containerID return the same IP.
- Allocated IPs range from `.2` to `.254` in the pod subnet.
- Corrupt entries (invalid IP strings) are detected and reported as errors.

#### CNI Plugin

- CNI commands are dispatched based on `CNI_COMMAND` env var: `ADD`, `DEL`, `CHECK`, `VERSION`.
- `VERSION` command does not require config or host IP.
- Plugin logs invocation context and errors.
- Writes errors in CNI-compliant JSON format to stdout/stderr.

#### External Command Usage

- Only one exec.Command call: `fanctl up` in `internal/fan/fanctl.go`.
- If `fanctl` is not found, error is explicit: "fanctl not found in PATH: install ubuntu-fan package".
- Output from `fanctl` is included in error messages for troubleshooting.

#### File Paths

- IPAM state is stored in `/var/lib/cni/fancni/` by default.
- Log file is `/var/log/fancni.log`.
- Helm chart is in `deploy/helm/fancni/`.

#### Environment Variables

- CNI plugin expects:
  - `CNI_COMMAND`
  - `CNI_CONTAINERID`
  - CNI config via stdin

#### Helm & Docker

- Helm chart values are in `deploy/helm/fancni/values.yaml`.
- Dockerfiles:
  - CNI plugin: `deploy/docker/Dockerfile.cni`
  - Init container: `deploy/docker/Dockerfile.init`
- Init scripts: `deploy/scripts/init-node.sh`, `deploy/scripts/install-cni.sh`

#### E2E Testing

- E2E tests use shell scripts, LXC VMs, and Kubernetes nodes.
- Test script (`tests/e2e/test-e2e.sh`) covers:
  - Multi-node setup
  - Container networking
  - Cross-node forwarding
  - Helm chart deployment
  - Validation of connectivity

#### Go Version

- Go module specifies `go 1.24.13`.
- Ensure Dockerfiles and CI use matching Go version.

## Summary Table

| Domain        | File/Package                    | Key Functions/Types         | Notes                                  |
|---------------|---------------------------------|-----------------------------|----------------------------------------|
| Fan           | `internal/fan/fan.go`           | `ComputeSubnet`, `ComputeGateway`, `ComputeBridgeName`, `ComputeUnderlayArg` | Pure functions, deterministic mapping  |
| Fanctl        | `internal/fan/fanctl.go`        | `EnsureBridge`              | Only exec.Command usage                |
| IPAM          | `internal/ipam/ipam.go`, `file_ipam.go` | `IPAM` interface, `FileIPAM` struct, `Allocate`, `Lookup`, `Free` | File-backed, flock, idempotent         |
| Config        | `internal/config/config.go`     | `Parse`, `NetConfig`        | Reads from stdin                       |
| CNI Plugin    | `internal/cni/plugin.go`        | `NewPlugin`, `HandleAdd`, `HandleDel`, `HandleCheck`, `HandleVersion` | Command dispatch                       |
| Logging       | `cmd/fancni/main.go`            | `log.SetOutput`             | `/var/log/fancni.log` or stderr        |
| Testing       | `<file>_test.go`                | `Test<Function>_<Scenario>` | Go `testing` package, helper functions |
| E2E           | `tests/e2e/test-e2e.sh`         | Shell script                | Multi-node, connectivity validation    |
| Helm          | `deploy/helm/fancni/`           | Chart.yaml, values.yaml     | Template, lint, deploy                 |
| Docker        | `deploy/docker/Dockerfile.*`    |                            | CNI/init containers                    |

## Quick Reference

- **Build:** `make build`
- **Test:** `make test`
- **E2E:** `make e2e`
- **Docker:** `make docker-build`
- **Helm:** `make helm-template`, `make helm-lint`
- **Log file:** `/var/log/fancni.log`
- **IPAM state:** `/var/lib/cni/fancni/ipam.json`
- **Fanctl error:** "fanctl not found in PATH: install ubuntu-fan package"
- **Overlay network:** Must be IPv4 CIDR (e.g., `"240.0.0.0/8"`)
- **Allocated IPs:** `.2`–`.254` in pod subnet
- **CNI commands:** `ADD`, `DEL`, `CHECK`, `VERSION`

---

**For new developers:**  
- Always wrap errors for context.
- Use pure functions for networking math.
- File-backed IPAM is concurrency-safe via flock.
- Logging defaults to `/var/log/fancni.log`, fallback to stderr.
- E2E tests are shell scripts, not Go tests.
- Only one exec.Command usage (`fanctl`), check PATH.
- Helm and Docker are first-class deployment artifacts.
- Follow naming conventions for clarity and consistency.
