# Task 002 Completion Report: Conditional FORWARD Rules

## Summary

- Wrapped the FORWARD ACCEPT iptables rules in `init-node.sh` with a conditional check on `NETWORK_POLICY_ENABLED != "true"`, so kube-router can manage the FORWARD chain unimpeded when network policy is enabled. The MASQUERADE rule remains unconditional.
- Added `NETWORK_POLICY_ENABLED` env var to the `node-init` init container in `daemonset.yaml`, sourced from `.Values.networkPolicy.enabled`.
- `make helm-lint` and `make helm-template` both pass; the rendered template shows `NETWORK_POLICY_ENABLED: "false"` by default.

## Files Changed

- `deploy/scripts/init-node.sh`
- `deploy/helm/fancni/templates/daemonset.yaml`

## Notable Tradeoffs / Risks

None. The change is minimal and backward-compatible — default value of `networkPolicy.enabled` is `false`, so existing deployments retain current behavior.
