# E2E Test Run Report — Full Rerun #2

**Date:** 2026-05-13  
**Command:** `make e2e`  
**Result:** ✅ PASS

---

## Phase 12 Summary Output

```
[e2e] === Phase 12: Summary ===
[e2e] IP assertion:   PASS=4  FAIL=0
[e2e] HTTP assertion: PASS=4   FAIL=0
[e2e] === E2E RESULT: PASS ===
```

---

## Phase-by-Phase Overview

| Phase | Description | Result |
|-------|-------------|--------|
| 1 | Creating LXC VMs (fancni-e2e-node-1, fancni-e2e-node-2) | ✅ |
| 2 | Waiting for cloud-init on both nodes | ✅ |
| 3 | Installing k8s snap (v1.35.3) on both nodes | ✅ |
| 4 | Bootstrapping fancni-e2e-node-1 (no CNI) | ✅ |
| 5 | Joining fancni-e2e-node-2 to the cluster | ✅ |
| 6 | Building fancni images on host (`fancni:latest`, `fancni-init:latest`) | ✅ (cached) |
| 7 | Transferring fancni, fancni-init, and nginx images into both VMs | ✅ |
| 8 | Installing fancni via Helm into kube-system | ✅ |
| 9 | Deploying nginx with 4 replicas | ✅ |
| 10 | Asserting pod IPs are in 240.0.0.0/8 | ✅ (4/4) |
| 11 | Asserting HTTP 200 from all pod IPs | ✅ (4/4) |
| 12 | Summary | ✅ PASS |

---

## Pod IP Details

| Pod IP | In 240.0.0.0/8 | HTTP 200 |
|--------|----------------|----------|
| 240.119.153.4 | PASS | PASS |
| 240.119.153.3 | PASS | PASS |
| 240.119.252.2 | PASS | PASS |
| 240.119.252.3 | PASS | PASS |

Two pods were scheduled on each node (`.153.x` and `.252.x` subnets), confirming per-node subnet allocation is working correctly.
