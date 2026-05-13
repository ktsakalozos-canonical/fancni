# Report: 004-fix-wait-ready

## Summary

- Replaced `k8s status --wait-ready --timeout 300s` with `wait_for 120 "k8s API responding on ${NODE1}" lxc exec "${NODE1}" -- sudo k8s kubectl get nodes` on line 105 of `tests/e2e/test-e2e.sh`.
- The old check failed because nodes are `NotReady` until a CNI is installed, which happens after bootstrap. The new check only waits for the API server to respond, which is the correct gate at this stage.

## Files Changed

- `tests/e2e/test-e2e.sh`

## Notable Tradeoffs or Risks

- None. The change is a direct substitution with no behavioral impact beyond removing the premature readiness gate.
