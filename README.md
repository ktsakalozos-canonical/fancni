# FanCNI

A lightweight Kubernetes CNI plugin that implements Ubuntu Fan networking, enabling flat, routable pod-to-pod connectivity across nodes without an overlay encapsulation overhead.

> ⚠️ **Alpha** — not yet production-ready.

## How It Works

Ubuntu Fan networking maps a large overlay address space (default `240.0.0.0/8`) onto the existing underlay network. Each node receives a `/16` slice of the overlay derived deterministically from its underlay IP: for a node with underlay address `A.B.C.D`, pods are assigned addresses in `240.A.B.0/24` (with the exact range governed by the configured prefix length). The kernel's built-in Fan tunnel driver forwards packets between nodes by encapsulating them in UDP and routing them over the underlay, so pods on different nodes can communicate without any additional overlay daemons or encapsulation libraries.

## Prerequisites

- Ubuntu nodes (Fan networking relies on Ubuntu kernel support)
- [Canonical Kubernetes](https://ubuntu.com/kubernetes) cluster (recommended)
- Helm 3
- Docker (for building images)

## Quick Start

### 1. Build container images

```bash
make docker-build        # builds fancni:latest and fancni-init:latest
```

### 2. Load images onto cluster nodes

```bash
docker save fancni:latest | sudo ctr -n k8s.io images import -
docker save fancni-init:latest | sudo ctr -n k8s.io images import -
```

### 3. Install with Helm

```bash
helm install fancni deploy/helm/fancni/
```

### 4. Verify

```bash
kubectl -n kube-system get pods -l app=fancni
# All pods should reach Running status within ~30 s
```

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for a full description of the components and design rationale.

## Configuration

See [docs/configuration.md](docs/configuration.md) for all available Helm values and configuration options.

## Connectivity Test

After installation, verify pod connectivity with the provided test manifests:

```bash
# Edit deploy/test/connectivity-test.yaml to set nodeSelector values
# that match two nodes in your cluster, then apply:
kubectl apply -f deploy/test/connectivity-test.yaml

kubectl get pods -o wide          # note the IPs assigned to each pod

kubectl exec fancni-test-1 -- ping -c3 <fancni-test-2-IP>
kubectl exec fancni-test-2 -- ping -c3 <fancni-test-1-IP>
kubectl exec fancni-test-1 -- ping -c3 8.8.8.8   # external connectivity
```

## Troubleshooting

**`fanctl: command not found`**  
The `fanctl` tool is not installed on the node. Ensure the `fancni-init` DaemonSet pod completed successfully on that node (`kubectl -n kube-system logs <fancni-init-pod>`).

**Bridge not created / fan device missing**  
Verify the node is running Ubuntu with a kernel that includes Fan support (`uname -r`; typically 4.4+). Check dmesg for fan-related errors.

**Pod stuck in `ContainerCreating`**  
The CNI plugin failed to configure the pod network. Inspect kubelet logs (`journalctl -u kubelet`) and the CNI log file (typically `/var/log/fancni.log` or the path configured in the CNI conf). Confirm the CNI binary exists at `/opt/cni/bin/fancni` on the node.

**Pods cannot reach other nodes' pods**  
Ensure all nodes are on the same L2 underlay network (Fan relies on direct underlay reachability). Check that the Fan tunnel interface (`fan-<overlay>`) exists on each node (`ip link show`).

## Development

See [docs/development.md](docs/development.md) for build instructions, test commands, and contribution guidelines.
