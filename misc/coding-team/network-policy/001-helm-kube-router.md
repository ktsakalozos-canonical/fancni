# Task: Add kube-router network policy support to fancni Helm chart

## Context

fancni is a CNI plugin deployed via a Helm chart at `deploy/helm/fancni/`. It currently has no network policy enforcement. We want to add kube-router (upstream image `cloudnativelabs/kube-router`) running in firewall-only mode as an opt-in second DaemonSet within the same chart.

## Objective

Add the following to the existing Helm chart, all conditional on `.Values.networkPolicy.enabled`:

1. **values.yaml** — add a `networkPolicy` section
2. **ClusterRole + ClusterRoleBinding** — kube-router needs read access to k8s API
3. **DaemonSet** — runs kube-router in firewall-only mode

## Scope

Files to create/modify under `deploy/helm/fancni/`:

- `values.yaml` — add `networkPolicy` block
- `templates/network-policy-clusterrole.yaml` (new)
- `templates/network-policy-daemonset.yaml` (new)

## values.yaml additions

```yaml
networkPolicy:
  enabled: false
  image:
    repository: cloudnativelabs/kube-router
    tag: v2.9.0
    pullPolicy: IfNotPresent
  iptablesSyncPeriod: "5m"
```

## ClusterRole permissions needed

kube-router firewall controller needs:
- `pods` — get, list, watch
- `namespaces` — get, list, watch
- `nodes` — get, list, watch
- `networkpolicies` (networking.k8s.io) — get, list, watch
- `services` — get, list, watch (needed for policy selectors)
- `endpoints` — get, list, watch
- `endpointslices` (discovery.k8s.io) — get, list, watch

Use the existing ServiceAccount (`{{ .Release.Name }}-fancni`) — no new SA needed.

## DaemonSet spec

- Name: `{{ .Release.Name }}-kube-router`
- Labels: `app: kube-router`
- hostNetwork: true
- serviceAccountName: same as fancni (`{{ .Release.Name }}-fancni`)
- Same tolerations as fancni (from `.Values.tolerations`)
- Single container:
  - image from `.Values.networkPolicy.image`
  - args: `--run-firewall=true`, `--run-router=false`, `--run-service-proxy=false`, `--run-loadbalancer=false`, `--iptables-sync-period={{ .Values.networkPolicy.iptablesSyncPeriod }}`
  - securityContext: privileged: true
  - volumeMounts: `/lib/modules` (read-only, for iptables/ipset kernel modules), `/run/xtables.lock` (for iptables lock coordination)
  - resources: requests cpu 100m, memory 64Mi (sensible defaults)
  - livenessProbe: httpGet /healthz port 20244, initialDelaySeconds 10, periodSeconds 10

## Constraints

- All new templates must be wrapped in `{{- if .Values.networkPolicy.enabled }}`
- Follow the existing chart's naming/label conventions
- Do NOT modify the existing daemonset.yaml or configmap.yaml
- ClusterRole (not Role) is required since network policies and pods span namespaces

## Non-goals

- Do not modify any Go source code
- Do not add a separate ServiceAccount for kube-router
- Do not add Helm tests
