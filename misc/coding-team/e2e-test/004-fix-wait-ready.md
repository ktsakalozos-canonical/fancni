# Task: Remove premature --wait-ready check

## Context

`k8s status --wait-ready` after bootstrap fails because nodes are `NotReady` until a CNI is installed — which is by design since we bootstrap without CNI. The bootstrap command itself already succeeded.

## Fix

In `tests/e2e/test-e2e.sh`, replace the `k8s status --wait-ready --timeout 300s` call (line 105) with a simpler check that the K8s API is responding:

```bash
wait_for 120 "k8s API responding on ${NODE1}" \
  lxc exec "${NODE1}" -- sudo k8s kubectl get nodes
```

This waits until `kubectl get nodes` succeeds (API server is up), without requiring nodes to be in `Ready` state. The nodes will become Ready after fancni is installed later.
