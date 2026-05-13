# Fancni Project Knowledge Reference

---

## Dependency List

**Go Modules (from `go.mod`):**
- Go version: `1.24.13`
- github.com/coreos/go-iptables `v0.8.0` (indirect)
- github.com/vishvananda/netlink `v1.3.1` (indirect)
- github.com/vishvananda/netns `v0.0.5` (indirect)
- golang.org/x/sys `v0.10.0` (indirect)

**System Dependencies:**
- `fanctl` binary (from `ubuntu-fan` package) must be installed on host for bridge setup.

---

## Build & Deploy Tooling

### Makefile Targets

- `build`: Compile fancni CNI binary to `_output/bin/fancni` from `./cmd/fancni/`.
- `test`: Run all Go tests with verbose output and no caching.
- `clean`: Remove `_output/` directory.
- `docker-build-cni`: Build Docker image for CNI plugin using `deploy/docker/Dockerfile.cni` (tag: `fancni:latest`).
- `docker-build-init`: Build Docker image for init container using `deploy/docker/Dockerfile.init` (tag: `fancni-init:latest`).
- `docker-build`: Build both CNI and init Docker images.
- `helm-template`: Render Helm chart templates for fancni.
- `helm-lint`: Lint Helm chart.
- `e2e`: Run end-to-end tests via `tests/e2e/test-e2e.sh`.

### Dockerfiles

- `deploy/docker/Dockerfile.cni`: Builds the CNI plugin image.
- `deploy/docker/Dockerfile.init`: Builds the init container image.

### Helm Chart

- Located at `deploy/helm/fancni/`
    - `Chart.yaml`, `values.yaml`, and templates for service account, configmap, daemonset.

### Scripts

- `deploy/scripts/install-cni.sh`: Installs the CNI plugin.
- `deploy/scripts/init-node.sh`: Initializes node for fancni.

---

## CI Workflow Inventory

### GitHub Actions Workflows (`.github/workflows/`):

- `nightly-dreams.yml`
    - Runs nightly, likely for project health, scalability, modularity, testing, etc. (see `misc/nightly-dreams/` for details).
    - May include validation, performance, security, automation, maintainability, observability, integration checks.

- `nightly-knowledge.yml`
    - Runs nightly, focuses on knowledge distillation and agent prompt updates.
    - Generates workflow file reports, prompt update reports, and includes prior knowledge.
    - Output files in `misc/coding-team/nightly-knowledge/`.

---

## Project Structure Map

### Top-level

- `cmd/fancni/`: Main CNI plugin entrypoint (`main.go`).
- `internal/`: Core implementation, subdivided:
    - `fan/`: Fan networking logic (`fan.go`, `fanctl.go`, tests).
    - `ipam/`: IP address management (`ipam.go`, `file_ipam.go`, tests).
    - `config/`: CNI config parsing (`config.go`, tests).
    - `netutil/`: Netlink helpers (`netutil.go`, tests).
    - `iptables/`: IPTables manipulation (`iptables.go`, tests).
    - `cni/`: Plugin orchestration logic (`plugin.go`, tests).
- `deploy/`: Deployment resources:
    - `docker/`: Dockerfiles for CNI and init containers.
    - `scripts/`: Shell scripts for install/init.
    - `helm/`: Helm chart for Kubernetes deployment.
    - `test/`: Example test manifests (`connectivity-test.yaml`).
- `tests/e2e/`: End-to-end test scripts (`test-e2e.sh`).
- `misc/`: Documentation, planning, and agent reports.
    - `nightly-dreams/`: Project health topics.
    - `coding-team/`: Knowledge and E2E test reports.
- `.opencode/agents/`: Agent prompt files for repo-scout, code-reviewer, developer, architect.

### Hotspots (recently changed/critical):

- `.github/workflows/nightly-knowledge.yml`: New knowledge distillation workflow.
- `misc/coding-team/nightly-knowledge/`: Knowledge reports and workflow file updates.
- `tests/e2e/test-e2e.sh`: Major additions for E2E testing (multi-node, LXC VMs, HTTP retries, Helm push, cross-node forwarding).
- `deploy/docker/Dockerfile.cni`, `deploy/scripts/init-node.sh`: Updated for E2E and Docker build fixes.
- `internal/fan/`, `internal/ipam/`: Core logic for fan networking and IPAM.
- `cmd/fancni/main.go`: Entrypoint, orchestrates config, IPAM, fan, plugin logic.

---

## Key Implementation Details

### CNI Plugin Entrypoint (`cmd/fancni/main.go`)
- Handles CNI commands: `ADD`, `DEL`, `CHECK`, `VERSION`.
- Reads config from stdin, detects host IP, computes pod subnet via fan networking.
- Uses file-backed IPAM for container IP allocation.
- Logs to `/var/log/fancni.log`.

### Fan Networking (`internal/fan/`)
- Deterministic mapping from underlay (node IP) to overlay (pod subnet).
- Pure functions for subnet/gateway/bridge name calculation.
- `fanctl` used for bridge setup (only exec call).

### IPAM (`internal/ipam/`)
- File-backed, concurrency-safe IP allocation.
- Allocates first free IP in pod subnet, idempotent per containerID.

### E2E Testing (`tests/e2e/test-e2e.sh`)
- Multi-node Kubernetes cluster setup (LXC VMs).
- Validates cross-node pod connectivity, HTTP retries, Helm chart deployment.
- Handles Docker image build/push, waits for node readiness, validates NGINX image pull.

---

## Actionable Guidance

- **Build:** Use `make build` for local binary, `make docker-build` for container images.
- **Deploy:** Use Helm chart (`deploy/helm/fancni/`) for Kubernetes, or scripts for manual install.
- **Test:** Use `make test` for unit tests, `make e2e` for full E2E validation.
- **CI:** Nightly workflows provide health and knowledge checks; review outputs in `misc/coding-team/nightly-knowledge/`.
- **Dependencies:** Ensure Go 1.24.13+ and `fanctl` (Ubuntu Fan) are installed on host.
- **Critical Files:** Monitor changes in `internal/fan/`, `internal/ipam/`, `cmd/fancni/main.go`, `.github/workflows/`, and `tests/e2e/`.

---

## Reference: Directory Map

```
cmd/fancni/                # Main CNI plugin
internal/fan/              # Fan networking logic
internal/ipam/             # IP address management
internal/config/           # CNI config parsing
internal/netutil/          # Netlink helpers
internal/iptables/         # IPTables logic
internal/cni/              # Plugin orchestration
deploy/docker/             # Dockerfiles
deploy/scripts/            # Install/init scripts
deploy/helm/fancni/        # Helm chart
deploy/test/               # Test manifests
tests/e2e/                 # E2E test scripts
misc/nightly-dreams/       # Project health docs
misc/coding-team/          # Knowledge/E2E reports
.opencode/agents/          # Agent prompts
.github/workflows/         # CI workflows
```

---

## CI Workflow Outputs

- **nightly-knowledge.yml**: Generates markdown reports on workflow files, prompt updates, prior knowledge inclusion, YAML indentation fixes.
    - Output: `misc/coding-team/nightly-knowledge/00*-*.md`

- **nightly-dreams.yml**: Topics include prioritization, integration, automation, scalability, observability, maintainability, modularity, testing, validation, performance, security.
    - Output: `misc/nightly-dreams/*.md`

---

## Most-Critical/Most-Changed Areas

- `tests/e2e/test-e2e.sh`: E2E test logic, multi-node setup, image pull, connectivity.
- `.github/workflows/nightly-knowledge.yml`: Knowledge distillation, agent prompt updates.
- `misc/coding-team/nightly-knowledge/`: Knowledge reports.
- `internal/fan/`, `internal/ipam/`: Core networking and IPAM logic.
- `cmd/fancni/main.go`: Entrypoint, orchestrates plugin operation.

---

## Quickstart Summary

1. **Build plugin:** `make build`
2. **Build containers:** `make docker-build`
3. **Deploy:** Use Helm chart or scripts.
4. **Run tests:** `make test` (unit), `make e2e` (integration)
5. **Check CI outputs:** See `misc/coding-team/nightly-knowledge/` for knowledge reports.

---

## Version Pinning

- Go: `1.24.13`
- All dependencies pinned via `go.mod` (see above).

---

## Agent Reference

- For repo-scout: Focus on `internal/fan/`, `internal/ipam/`, `.github/workflows/`, `tests/e2e/`, and `misc/coding-team/nightly-knowledge/` for most critical/changed logic.
- For build/deploy: Use Makefile and Dockerfiles; Helm chart for Kubernetes.
- For CI: Nightly workflows generate actionable reports and project health documentation.

---
