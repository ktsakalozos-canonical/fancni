# Configuration Reference

The following Helm chart values can be used to configure fancni. They are defined in `deploy/helm/fancni/values.yaml`.

| Key | Default | Description |
|-----|---------|-------------|
| `overlay.network` | `240.0.0.0/8` | Fan overlay address space |
| `underlay.prefix` | `16` | Underlay subnet prefix length used to derive per-node pod CIDRs |
| `images.cni.repository` | `fancni` | CNI container image name |
| `images.cni.tag` | `latest` | CNI container image tag |
| `images.cni.pullPolicy` | `IfNotPresent` | Image pull policy for the CNI container |
| `images.init.repository` | `fancni-init` | Init container image name |
| `images.init.tag` | `latest` | Init container image tag |
| `images.init.pullPolicy` | `IfNotPresent` | Image pull policy for the init container |
| `cni.name` | `fancni` | CNI plugin name |
| `cni.confFileName` | `10-fancni.conf` | CNI config file name written to `/etc/cni/net.d/` |
| `cni.version` | `1.0.0` | CNI specification version used in the config file |
| `ipam.dataDir` | `/var/lib/cni/fancni` | Directory for IPAM state files |
| `tolerations` | `[{operator: Exists, effect: NoSchedule}, {operator: Exists, effect: NoExecute}]` | DaemonSet tolerations (ensures pods run on all nodes including tainted ones) |
| `networkPolicy.enabled` | `true` | Enable network policy enforcement via kube-router |
| `networkPolicy.image.repository` | `cloudnativelabs/kube-router` | Network policy controller image |
| `networkPolicy.image.tag` | `v2.9.0` | Network policy controller image tag |
| `networkPolicy.image.pullPolicy` | `IfNotPresent` | Image pull policy for network policy controller |
| `networkPolicy.iptablesSyncPeriod` | `5m` | How often iptables rules are synced |
