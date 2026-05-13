# Task 006: Helm Chart — Completion Report

## Summary

- Created Helm chart under `deploy/helm/fancni/` with `Chart.yaml`, `values.yaml`, and three templates (`configmap.yaml`, `serviceaccount.yaml`, `daemonset.yaml`).
- DaemonSet uses `hostNetwork: true`, `hostPID: true`, one privileged init container (`node-init`), and one main container (`install-cni`) that reads the CNI config from the ConfigMap via `CNI_CONF` env var.
- Tolerations cover all taints (`NoSchedule` and `NoExecute` with `operator: Exists`) to ensure scheduling on all nodes including control plane.
- Added `helm-template` and `helm-lint` targets to the Makefile; `helm lint` passes with 0 failures.

## Files Changed

- `deploy/helm/fancni/Chart.yaml` (new)
- `deploy/helm/fancni/values.yaml` (new)
- `deploy/helm/fancni/templates/configmap.yaml` (new)
- `deploy/helm/fancni/templates/serviceaccount.yaml` (new)
- `deploy/helm/fancni/templates/daemonset.yaml` (new)
- `Makefile` (added `helm-template` and `helm-lint` targets)

## Notable Tradeoffs or Risks

- None. The chart is intentionally minimal — no RBAC, no PSP, no Services — per the task brief constraints.
