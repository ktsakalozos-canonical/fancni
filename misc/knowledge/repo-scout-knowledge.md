## Dependency List

**Go Modules (from `go.mod`):**
- Go version: `1.24.13`
- `github.com/coreos/go-iptables` `v0.8.0` (indirect)
- `github.com/vishvananda/netlink` `v1.3.1` (indirect)
- `github.com/vishvananda/netns` `v0.0.5` (indirect)
- `golang.org/x/sys` `v0.10.0` (indirect)

**System Dependencies:**
- `fanctl` binary (from `ubuntu-fan` package) must be installed on host for bridge setup.

---

## Build & Deploy Tooling

### Makefile Targets

- `build`: Compile `fancni` CNI binary to `_output/bin/fancni` from `./cmd/fancni/`.
- `test`: Run all Go tests with verbose output and no caching.
- `clean`: Remove `_output/` directory.
- `rock-build`: Package the application using Rockcraft.
- `helm-template`: Render Helm chart templates for `fancni`.
- `helm-lint`: Lint Helm chart.
- `e2e`: Run end-to-end tests via `tests/e2e/test-e2e.sh`.

### Dockerfiles

- **Removed Docker references**: All Docker-related files have been replaced with Rockcraft configurations.

### Helm Chart

- Located at `deploy/helm/fancni/`
    - `Chart.yaml`, `values.yaml`, and templates for service account, configmap, daemonset.

### Scripts

- `deploy/scripts/install-cni.sh`: Installs the CNI plugin.
- `deploy/scripts/init-node.sh`: Initializes node for `fancni`.

---

## CI Workflow Inventory

### GitHub Actions Workflows (`.github/workflows/`):

- `nightly-dreams.yml`
    - Runs nightly, focusing on project health, scalability, modularity, testing, and more.
  
- `release-latest.yml`
    - Manages release process for the latest version of the project.
  
- `ci-tests.yml`
    - Executes continuous integration tests, including unit tests and end-to-end tests.
  
- `nightly-knowledge.yml`
    - Runs nightly, focuses on knowledge distillation and agent prompt updates.
    - Generates workflow file reports, prompt update reports, and includes prior knowledge.

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
    - `helm/`: Helm chart for Kubernetes deployment.
    - `scripts/`: Shell scripts for install/init.
- `tests/e2e/`: End-to-end test scripts (`test-e2e.sh`).
- `misc/`: Documentation, planning, and agent reports.
    - `nightly-dreams/`: Project health topics.
    - `coding-team/`: Knowledge and E2E test reports.
- `.opencode/agents/`: Agent prompt files for repo-scout, code-reviewer, developer, architect.

### Hotspots (recently changed/critical):

- `.github/workflows/nightly-knowledge.yml`: New knowledge distillation workflow.
- `tests/e2e/test-e2e.sh`: Major additions for E2E testing (multi-node, LXC VMs, HTTP retries, Helm push, cross-node forwarding).
- `internal/fan/`, `internal/ipam/`: Core logic for fan networking and IPAM.
- `cmd/fancni/main.go`: Entrypoint, orchestrates config, IPAM, fan, plugin logic.

--- 

## Key Implementation Details

- The project has transitioned from Docker to Rockcraft for packaging and deployment.
- CI workflows have been updated to include new tests and knowledge distillation processes.
- The `Makefile` has been streamlined to focus on Go build and Helm-related tasks.
