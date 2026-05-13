# Task: E2E test script and Makefile target

## Context

fancni is a Kubernetes CNI plugin using Ubuntu Fan networking. There are no automated e2e tests. The project uses Canonical K8s (`k8s` snap), Go, Helm, and Docker. See `README.md` and `deploy/` for deployment details.

## Objective

Create a bash e2e test script at `tests/e2e/test-e2e.sh` and a `make e2e` Makefile target.

## What the script does (in order)

1. **Create two LXC VMs** named `fancni-e2e-node-1` and `fancni-e2e-node-2` using `lxc launch ubuntu:24.04 --vm`. Give them reasonable resources (2 CPU, 4GB+ RAM, 20GB disk).

2. **Wait for VMs to be ready** — wait until `lxc exec <vm> -- cloud-init status --wait` succeeds on both.

3. **Install the `k8s` snap** on both VMs: `lxc exec <vm> -- sudo snap install k8s --classic --channel=1.35-classic/stable`

4. **Bootstrap node-1 without CNI** — write a bootstrap config that disables the default network, then bootstrap:
   ```
   cat <<EOF > bootstrap-config.yaml
   cluster-config:
     network:
       enabled: false
   EOF
   sudo k8s bootstrap --file bootstrap-config.yaml
   ```
   Wait for `sudo k8s status --wait-ready` to succeed.

5. **Join node-2** — on node-1 run `sudo k8s get-join-token fancni-e2e-node-2 --worker`, then on node-2 run `sudo k8s join-cluster <token>`. Wait until `sudo k8s kubectl get nodes` shows both nodes.

6. **Build fancni images on the host** — run `make docker-build` from the repo root.

7. **Transfer images into both VMs** — for each VM and each image (`fancni:latest`, `fancni-init:latest`):
   ```
   docker save <image> | lxc exec <vm> -- /snap/k8s/current/bin/ctr \
     --address /var/snap/k8s/common/run/containerd.sock \
     -n k8s.io images import -
   ```

8. **Install fancni via Helm** — on node-1:
   ```
   sudo k8s helm install fancni <chart-path> --namespace kube-system
   ```
   The chart lives on the host at `deploy/helm/fancni/`. Push the chart directory into the VM first (e.g. `lxc file push --recursive deploy/helm/fancni/ <vm>/tmp/fancni-chart/`), then helm install from `/tmp/fancni-chart/fancni`.
   Wait until the fancni DaemonSet pods are Running on both nodes.

9. **Deploy nginx** — create a deployment with 4 replicas:
   ```
   sudo k8s kubectl create deployment nginx-e2e --image=nginx:latest --replicas=4
   ```
   Wait until all 4 pods are Running (timeout ~120s).

10. **Assert pods got fancni IPs** — get pod IPs via `kubectl get pods -o wide`. Verify each IP starts with `240.` (the overlay network from `values.yaml`).

11. **Assert HTTP responses** — for each pod IP, run `curl -s -o /dev/null -w '%{http_code}' http://<pod-ip>` from inside one of the VMs (since overlay IPs are routable within cluster nodes). Expect HTTP 200.

12. **Print summary** — log pass/fail for each assertion.

## Cleanup

- Use a `trap` to clean up on exit: `lxc delete --force fancni-e2e-node-1 fancni-e2e-node-2`
- Accept a `--no-cleanup` flag (as `$1` or env var `FANCNI_E2E_NO_CLEANUP=1`) that skips the cleanup trap for debugging.

## Makefile

Add to the existing `Makefile`:
```
e2e:
	bash tests/e2e/test-e2e.sh
```
Add `e2e` to the `.PHONY` list.

## Non-goals
- No CI integration
- No multi-arch
- No DNS/service/ingress testing
- No performance testing

## Constraints / Caveats
- The `k8s` snap does NOT have a `ctr` subcommand. Use the binary directly: `/snap/k8s/current/bin/ctr --address /var/snap/k8s/common/run/containerd.sock -n k8s.io`
- `k8s helm` and `k8s kubectl` are the correct commands (not standalone helm/kubectl)
- Canonical K8s bootstrap config to disable CNI uses the YAML structure shown above
- The script should `set -euo pipefail` and use clear `echo` statements to log each phase
- Use reasonable timeouts and retries for async operations (snap install, bootstrap, pod readiness). A helper `wait_for` function that retries a command with a timeout is recommended.
- Make the script executable (`chmod +x`)
