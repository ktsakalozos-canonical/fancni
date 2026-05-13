# Task: Fix containerd socket path

## Context

The containerd socket in Canonical K8s is at `/run/containerd/containerd.sock`, NOT `/var/snap/k8s/common/run/containerd.sock`. Verified by bootstrapping the k8s snap and checking the actual socket location.

## Fix

In `tests/e2e/test-e2e.sh`, change the `CTR_SOCK` variable from:
```
CTR_SOCK="/var/snap/k8s/common/run/containerd.sock"
```
to:
```
CTR_SOCK="/run/containerd/containerd.sock"
```

That's the only change needed.
