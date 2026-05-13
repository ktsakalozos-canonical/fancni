# Task 007: README and Test Manifests

## Context

All code, containers, and Helm chart are complete. This task adds a README and test manifests for validating pod connectivity.

## Objective

1. Write a `README.md` at the repo root with quickstart, architecture overview, and troubleshooting
2. Create test Kubernetes manifests for verifying pod connectivity

## Scope

### README.md

Structure:
1. **Title + one-line description**
2. **Prerequisites**: Ubuntu nodes, Kubernetes cluster (kubeadm, k3s, etc.), Helm 3, Docker
3. **Quick Start**: build images, helm install, verify
4. **Architecture**: brief description of components (CNI binary, init container, Helm chart), link to `misc/plan.md` for details
5. **Configuration**: table of Helm values
6. **How It Works**: one paragraph on fan networking address mapping
7. **Troubleshooting**: common issues (fanctl not found, bridge not created, pod stuck in ContainerCreating)
8. **Development**: how to build, test, lint locally

Keep it concise — no more than ~150 lines total.

### Test manifests — deploy/test/

**deploy/test/connectivity-test.yaml**

A simple test deployment:
- 2 pods with `nodeSelector` to land on different nodes (if available)
- Each pod runs `busybox` with `sleep infinity`
- After applying, user can exec into one pod and ping the other's IP

```yaml
# Pod 1
apiVersion: v1
kind: Pod
metadata:
  name: fancni-test-1
  labels:
    app: fancni-test
spec:
  containers:
    - name: test
      image: busybox:1.36
      command: ["sleep", "infinity"]
---
# Pod 2
apiVersion: v1
kind: Pod
metadata:
  name: fancni-test-2
  labels:
    app: fancni-test
spec:
  containers:
    - name: test
      image: busybox:1.36
      command: ["sleep", "infinity"]
```

Instructions in the README for testing:
```bash
kubectl apply -f deploy/test/connectivity-test.yaml
kubectl get pods -o wide  # Note the IPs
kubectl exec fancni-test-1 -- ping -c3 <fancni-test-2-IP>
kubectl exec fancni-test-2 -- ping -c3 <fancni-test-1-IP>
kubectl exec fancni-test-1 -- ping -c3 8.8.8.8  # External connectivity
```

## Non-goals

- No automated CI pipeline
- No integration test framework
- No GitHub Actions / workflow files

## Constraints

- README must be accurate to the actual project structure
- Test manifests must be valid Kubernetes YAML
