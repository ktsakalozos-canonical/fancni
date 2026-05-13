# E2E Full Rerun Report

**Date:** 2026-05-13  
**Command:** `make e2e`  
**Verdict:** ❌ FAIL

---

## Summary

The e2e test ran from scratch (fresh LXC VMs) through Phase 8 successfully, but timed out in Phase 9 waiting for 4 nginx pods to reach `Running` state within 120 seconds. The VMs were cleaned up automatically before pod state could be inspected.

---

## Phase-by-Phase Results

| Phase | Description | Result |
|-------|-------------|--------|
| 1 | Creating LXC VMs (fancni-e2e-node-1, fancni-e2e-node-2) | ✅ PASS |
| 2 | Waiting for cloud-init | ✅ PASS |
| 3 | Installing k8s snap (1.35-classic/stable) | ✅ PASS |
| 4 | Bootstrapping node-1 without CNI | ✅ PASS |
| 5 | Joining node-2 to cluster | ✅ PASS |
| 6 | Building fancni images on host | ✅ PASS |
| 7 | Transferring images into VMs | ✅ PASS |
| 8 | Installing fancni via Helm + waiting for DaemonSet ready | ✅ PASS |
| 9 | Deploying nginx (4 replicas) — **TIMEOUT** after 120s | ❌ FAIL |
| 10–12 | IP/HTTP assertions, Summary | Not reached |

---

## Failure Detail

```
[e2e] === Phase 9: Deploying nginx (4 replicas) ===
deployment.apps/nginx-e2e created
[e2e] Waiting for: 4 nginx pods Running (timeout 120s)
[e2e] TIMEOUT waiting for: 4 nginx pods Running
[e2e] Cleaning up VMs...
make: *** [Makefile:27: e2e] Error 1
```

The VMs were deleted by the cleanup trap before any pod diagnostics could be collected.

---

## Root Cause Analysis

The nginx pods failed to reach `Running` within 120 seconds. Likely causes (in order of probability):

1. **Image pull latency** — `nginx:latest` must be pulled from Docker Hub inside each LXC VM. On a cold cache, this can take >120 seconds depending on network speed. The fancni DaemonSet pods use `imagePullPolicy: Never` (pre-loaded images), but nginx uses `latest` which defaults to `IfNotPresent` or triggers a remote pull.

2. **CNI not fully ready for new pods** — The fancni DaemonSet was confirmed Running, but there may be a race where node networking is not yet fully operational for newly scheduled pods (routes/iptables not yet applied) causing pods to stay in `ContainerCreating`.

3. **120s timeout is too short** — The wait for nginx is only 120 seconds. If image pull takes 60–90s, remaining time for pod startup may be insufficient.

---

## Full Test Output

```
[e2e] === Phase 1: Creating LXC VMs ===
Launching fancni-e2e-node-1
[... image unpack progress ...]
Launching fancni-e2e-node-2
[... image unpack progress ...]
[e2e] === Phase 2: Waiting for cloud-init ===
[e2e] Waiting for: cloud-init fancni-e2e-node-1 (timeout 300s)
[e2e] Ready: cloud-init fancni-e2e-node-1
[e2e] Waiting for: cloud-init fancni-e2e-node-2 (timeout 300s)
[e2e] Ready: cloud-init fancni-e2e-node-2
[e2e] === Phase 3: Installing k8s snap ===
k8s (1.35-classic/stable) v1.35.3 from Canonical** installed
k8s (1.35-classic/stable) v1.35.3 from Canonical** installed
[e2e] === Phase 4: Bootstrapping fancni-e2e-node-1 (no CNI) ===
Bootstrapped a new Kubernetes cluster with node address "10.248.119.220:6400".
The node will be 'Ready' to host workloads after the CNI is deployed successfully.
[e2e] Waiting for: k8s API responding on fancni-e2e-node-1 (timeout 120s)
[e2e] Ready: k8s API responding on fancni-e2e-node-1
[e2e] === Phase 5: Joining fancni-e2e-node-2 to the cluster ===
[e2e] Join token obtained
Cluster services have started on "fancni-e2e-node-2".
[e2e] Waiting for: both nodes visible (timeout 300s)
[e2e] Ready: both nodes visible
[e2e] === Phase 6: Building fancni images on host ===
docker build -t fancni:latest -f deploy/docker/Dockerfile.cni .   [CACHED - success]
docker build -t fancni-init:latest -f deploy/docker/Dockerfile.init .  [CACHED - success]
[e2e] === Phase 7: Transferring images into VMs ===
[e2e]   Importing fancni:latest into fancni-e2e-node-1...   [success]
[e2e]   Importing fancni:latest into fancni-e2e-node-2...   [success]
[e2e]   Importing fancni-init:latest into fancni-e2e-node-1...  [success]
[e2e]   Importing fancni-init:latest into fancni-e2e-node-2...  [success]
[e2e] === Phase 8: Installing fancni via Helm ===
NAME: fancni
LAST DEPLOYED: Wed May 13 10:55:43 2026
NAMESPACE: kube-system
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
[e2e] Waiting for fancni DaemonSet pods to be Running...
[e2e] Waiting for: fancni DaemonSet ready (timeout 300s)
[e2e] Ready: fancni DaemonSet ready
[e2e] === Phase 9: Deploying nginx (4 replicas) ===
deployment.apps/nginx-e2e created
[e2e] Waiting for: 4 nginx pods Running (timeout 120s)
[e2e] TIMEOUT waiting for: 4 nginx pods Running
[e2e] Cleaning up VMs...
make: *** [Makefile:27: e2e] Error 1
```

---

## Recommendations

1. **Increase nginx wait timeout** from 120s to 300s in `test-e2e.sh` (line 160) to tolerate image pull latency.
2. **Add diagnostic output on timeout** — before cleaning up, dump `kubectl get pods -A -o wide` and `kubectl describe pods -l app=nginx-e2e` so failures are diagnosable.
3. **Consider pre-loading nginx image** into the VMs alongside fancni images (Phase 7) to eliminate pull latency from the critical path.
4. Optionally, use `--no-cleanup` flag (`FANCNI_E2E_NO_CLEANUP=1`) on next rerun to inspect cluster state after failure.
