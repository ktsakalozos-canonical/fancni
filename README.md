# fancni

A lightweight Kubernetes CNI plugin that implements Ubuntu Fan networking, enabling flat, routable pod-to-pod connectivity across nodes without an overlay encapsulation overhead.

> ⚠️ **Alpha** — not yet production-ready.

## How It Works

Ubuntu Fan networking maps a large overlay address space (default `240.0.0.0/8`) onto the existing underlay network. Each node receives a `/16` slice of the overlay derived deterministically from its underlay IP: for a node with underlay address `A.B.C.D`, pods are assigned addresses in `240.A.B.0/24` (with the exact range governed by the configured prefix length). The kernel's built-in Fan tunnel driver forwards packets between nodes by encapsulating them in UDP and routing them over the underlay, so pods on different nodes can communicate without any additional overlay daemons or encapsulation libraries.

See [ARCHITECTURE.md](ARCHITECTURE.md) for a full description of the components and design rationale.

## Prerequisites

- Ubuntu nodes (Fan networking relies on Ubuntu kernel support)
- [Canonical Kubernetes](https://ubuntu.com/kubernetes) cluster (recommended)
- Helm 3

## Quick Start

### 1. Deploy Canonical K8s without the network


```bash
cat <<EOF > bootstrap-config.yaml
cluster-config:
  network:
    enabled: false
  dns:
    enabled: true
EOF
sudo snap install k8s --classic --channel=1.35-classic/stable
sudo k8s bootstrap --file bootstrap-config.yaml
sudo k8s kubectl config view --raw > ~/.kube/config
```

### 2. Install with Helm

```bash
helm install fancni deploy/helm/fancni/
```

See [docs/configuration.md](docs/configuration.md) for all available Helm values and configuration options.


### 3. Verify

All pods should reach Running status within ~30s

```bash
sudo k8s kubectl get pods -l app=fancni
```


## Troubleshooting

**`fanctl: command not found`**  
The `fanctl` tool is not installed on the node. Ensure the `fancni-init` DaemonSet pod completed successfully on that node (`kubectl logs <fancni-init-pod>`).

**Bridge not created / fan device missing**  
Verify the node is running Ubuntu with a kernel that includes Fan support (`uname -r`; typically 4.4+). Check dmesg for fan-related errors.

**Pod stuck in `ContainerCreating`**  
The CNI plugin failed to configure the pod network. Inspect kubelet logs (`journalctl -u kubelet`) and the CNI log file (typically `/var/log/fancni.log` or the path configured in the CNI conf). Confirm the CNI binary exists at `/opt/cni/bin/fancni` on the node.

**Pods cannot reach other nodes' pods**  
Ensure all nodes are on the same L2 underlay network (Fan relies on direct underlay reachability). Check that the Fan tunnel interface (`fan-<overlay>`) exists on each node (`ip link show`).

## Development

See [docs/development.md](docs/development.md) for build instructions, test commands, and contribution guidelines.
