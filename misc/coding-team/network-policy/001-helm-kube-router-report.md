# Task Report: 001-helm-kube-router

## Summary

- Added `networkPolicy` block to `values.yaml` with `enabled: false` default, image config (`cloudnativelabs/kube-router:v2.9.0`), and `iptablesSyncPeriod`.
- Created `templates/network-policy-clusterrole.yaml`: ClusterRole granting read access to pods, namespaces, nodes, services, endpoints, networkpolicies, and endpointslices, plus a ClusterRoleBinding to the existing `{{ .Release.Name }}-fancni` ServiceAccount — both conditional on `networkPolicy.enabled`.
- Created `templates/network-policy-daemonset.yaml`: DaemonSet running kube-router in firewall-only mode (`--run-firewall=true`, all other modes disabled), with `hostNetwork`, privileged security context, `/lib/modules` (read-only) and `/run/xtables.lock` volume mounts, resource requests, and a liveness probe — conditional on `networkPolicy.enabled`.

## Files Changed

- `deploy/helm/fancni/values.yaml` — added `networkPolicy` section
- `deploy/helm/fancni/templates/network-policy-clusterrole.yaml` — new file
- `deploy/helm/fancni/templates/network-policy-daemonset.yaml` — new file

## Validation

`make helm-lint` and `make helm-template` both pass (0 chart failures). Confirmed new resources are absent by default and render correctly with `--set networkPolicy.enabled=true`.

## Notable Tradeoffs / Risks

- None. Existing daemonset.yaml and configmap.yaml were not modified. No new ServiceAccount created per spec.
