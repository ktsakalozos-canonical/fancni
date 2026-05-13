# Task: Update docs/configuration.md with missing Helm values

## Context
The `docs/configuration.md` was extracted from the old README which was incomplete. Comparing against `deploy/helm/fancni/values.yaml` reveals several missing keys.

## Objective
Update `docs/configuration.md` to include ALL values from `values.yaml`.

## Missing keys to add

| Key | Default | Description |
|-----|---------|-------------|
| `images.cni.pullPolicy` | `IfNotPresent` | Image pull policy for the CNI container |
| `images.init.pullPolicy` | `IfNotPresent` | Image pull policy for the init container |
| `cni.version` | `1.0.0` | CNI specification version used in the config file |
| `tolerations` | `[{operator: Exists, effect: NoSchedule}, {operator: Exists, effect: NoExecute}]` | DaemonSet tolerations (ensures pods run on all nodes including tainted ones) |
| `networkPolicy.enabled` | `true` | Enable network policy enforcement via kube-router |
| `networkPolicy.image.repository` | `cloudnativelabs/kube-router` | Network policy controller image |
| `networkPolicy.image.tag` | `v2.9.0` | Network policy controller image tag |
| `networkPolicy.image.pullPolicy` | `IfNotPresent` | Image pull policy for network policy controller |
| `networkPolicy.iptablesSyncPeriod` | `5m` | How often iptables rules are synced |

## Constraints
- Keep the same table format as existing entries
- Group logically (existing groups: overlay, underlay, images, cni, ipam — add tolerations, networkPolicy)
- Don't change existing descriptions

## Non-goals
- Don't modify any other files
