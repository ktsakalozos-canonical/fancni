# Task 003: E2E Network Policy Tests — Completion Report

## Summary

- **Flipped default**: `networkPolicy.enabled` changed from `false` to `true` in `values.yaml` so network policy enforcement is on by default.
- **Extended e2e test**: Added kube-router image to the Phase 7 transfer loop; added Phase 12 (kube-router readiness), Phase 13 (deny-all blocks traffic), Phase 14 (allow policy restores traffic); renumbered old Phase 12 (Summary) to Phase 15. All new phases follow existing `log`/`wait_for`/`lxc_exec` conventions and include retry loops for kube-router sync latency.
- **Updated plan.md**: Removed "No network policies: Out of scope" from Constraints and "Network policy enforcement" from Out of Scope; added a new `## Network Policy Support` section documenting the kube-router integration and default configuration.

## Files Changed

- `deploy/helm/fancni/values.yaml`
- `tests/e2e/test-e2e.sh`
- `misc/plan.md`

## Validation

`make helm-lint` passes (1 chart linted, 0 failed).

## Notable Tradeoffs / Risks

- The deny-phase uses a 15s hard sleep + 5-retry loop (10s apart) to handle kube-router's sync latency. The configured `iptablesSyncPeriod` is 5m, but new rules are also applied on policy-watch events, so 15s is typically sufficient. If the cluster is under load, the retry loop provides an additional ~50s buffer before declaring failure.
- The test creates pods via `kubectl run` (single pod, not a Deployment) in `netpol-test` namespace for simplicity; this matches the task brief.
