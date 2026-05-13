# Report: Quick Validation of Cross-Node Fix

## Verdict: PASS ✓

Cross-node pod connectivity is now working after identifying and fixing a secondary bug: the init script was using `iptables` (which resolves to `iptables-nft` on Ubuntu 24.04) instead of `iptables-legacy`, which is what kube-proxy uses. The fix required updating `deploy/scripts/init-node.sh` to use `iptables-legacy`.

---

## Steps Performed

### Step 1: Rebuild init image

```
docker build -t fancni-init:latest -f deploy/docker/Dockerfile.init .
```

**Note:** Build context had to be the workspace root (`.`), not `deploy/`, because the Dockerfile references `deploy/scripts/init-node.sh` as a path relative to the build context.

Build output: success, new image SHA `e94b84eb...` (first attempt with wrong context) → `0cebbf46...` (correct context).

### Step 2: Re-import image into both VMs

```
=== Importing into fancni-e2e-node-1 ===
docker.io/library/fancni init:latest    saved
application/vnd.oci.image.index.v1+json sha256:1eddc20071638edfe3b0fac9d57076206b9485090a7900455eecacfd628de5ac

=== Importing into fancni-e2e-node-2 ===
docker.io/library/fancni init:latest    saved
application/vnd.oci.image.index.v1+json sha256:1eddc20071638edfe3b0fac9d57076206b9485090a7900455eecacfd628de5ac
```

### Step 3: Delete fancni pods

```
pod "fancni-fancni-ktwhx" deleted from kube-system namespace
pod "fancni-fancni-ttnm4" deleted from kube-system namespace
```

### Step 4: Wait for pods to be Running

```
NAME                  READY   STATUS    RESTARTS   AGE
fancni-fancni-n7qbf   1/1     Running   0          39s
fancni-fancni-wc89b   1/1     Running   0          39s
```

### Step 5: Initial iptables verification (PROBLEM FOUND)

After the first pod restart, iptables rules appeared correct in `iptables-nft`:

```
=== Node 1 iptables (nft) FORWARD ===
Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     0    --  240.119.94.0/24      0.0.0.0/0
2    ACCEPT     0    --  0.0.0.0/0            240.119.94.0/24

=== Node 2 iptables (nft) FORWARD ===
Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     0    --  240.119.117.0/24     0.0.0.0/0
2    ACCEPT     0    --  0.0.0.0/0            240.119.117.0/24
```

But connectivity tests showed 000 (timeout) for cross-node pods:

```
Testing 240.119.117.4... 000    ← node-2 pod, FAIL
Testing 240.119.94.3... 200     ← node-1 pod, ok
Testing 240.119.117.3... 000    ← node-2 pod, FAIL
Testing 240.119.94.2... 200     ← node-1 pod, ok
```

**Root cause diagnosis:**

```
=== iptables-legacy FORWARD (used by kube-proxy) ===
Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination
1    KUBE-PROXY-FIREWALL  0  --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
2    KUBE-FORWARD         0  --  0.0.0.0/0  0.0.0.0/0
3    KUBE-SERVICES        0  --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
4    KUBE-EXTERNAL-SERVICES  0  --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
```

The system has **two independent iptables backends**: `iptables-legacy` (used by kube-proxy) and `iptables-nft`. The init script was writing to `iptables-nft`, but `iptables-legacy` was the active backend used by the kernel for packet forwarding decisions. Our ACCEPT rules were invisible to the forwarding path.

Tcpdump confirmed: ICMP replies from node-2 pods arrived at `fan-240 In` but were never forwarded out via `ftun0` because the legacy FORWARD chain had no ACCEPT rule for the pod subnet.

**Fix applied to `deploy/scripts/init-node.sh`:**
Changed `iptables` → `iptables-legacy` in the FORWARD ACCEPT rules and MASQUERADE rule:

```bash
# Before
nsenter --target 1 --net -- iptables -C FORWARD $RULE 2>/dev/null
nsenter --target 1 --net -- iptables -I FORWARD 1 $RULE
nsenter --target 1 --net -- iptables -t nat -C POSTROUTING $MASQ_RULE
nsenter --target 1 --net -- iptables -t nat -A POSTROUTING $MASQ_RULE

# After
nsenter --target 1 --net -- iptables-legacy -C FORWARD $RULE 2>/dev/null
nsenter --target 1 --net -- iptables-legacy -I FORWARD 1 $RULE
nsenter --target 1 --net -- iptables-legacy -t nat -C POSTROUTING $MASQ_RULE
nsenter --target 1 --net -- iptables-legacy -t nat -A POSTROUTING $MASQ_RULE
```

Image rebuilt and re-imported (SHA `f804f3b1...`), pods restarted again.

---

### Step 5 (final): iptables-legacy FORWARD chain verification

```
=== Node 1 iptables-legacy FORWARD ===
Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     0    --  0.0.0.0/0            240.119.94.0/24
2    ACCEPT     0    --  240.119.94.0/24      0.0.0.0/0
3    KUBE-PROXY-FIREWALL  0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
4    KUBE-FORWARD         0    --  0.0.0.0/0  0.0.0.0/0
5    KUBE-SERVICES        0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
6    KUBE-EXTERNAL-SERVICES  0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW

=== Node 2 iptables-legacy FORWARD ===
Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination
1    ACCEPT     0    --  0.0.0.0/0            240.119.117.0/24
2    ACCEPT     0    --  240.119.117.0/24     0.0.0.0/0
3    KUBE-PROXY-FIREWALL  0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
4    KUBE-FORWARD         0    --  0.0.0.0/0  0.0.0.0/0
5    KUBE-SERVICES        0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
6    KUBE-EXTERNAL-SERVICES  0    --  0.0.0.0/0  0.0.0.0/0  ctstate NEW
```

✓ ACCEPT rules for pod subnet at positions 1 and 2 on both nodes.

### Step 6: sysctl verification

```
=== Node 1 sysctls ===
net.ipv4.conf.fan-240.send_redirects = 0
net.ipv4.conf.fan-240.rp_filter = 0
net.ipv4.conf.all.rp_filter = 0

=== Node 2 sysctls ===
net.ipv4.conf.fan-240.send_redirects = 0
net.ipv4.conf.fan-240.rp_filter = 0
net.ipv4.conf.all.rp_filter = 0
```

✓ All sysctls are 0 as expected.

### Step 7: Cross-node connectivity test

```
=== Pod placement ===
NAME                         READY   STATUS    NODE
nginx-e2e-65999c5895-8z72v   1/1     Running   fancni-e2e-node-2   IP: 240.119.117.4
nginx-e2e-65999c5895-bd527   1/1     Running   fancni-e2e-node-1   IP: 240.119.94.3
nginx-e2e-65999c5895-mqcpw   1/1     Running   fancni-e2e-node-2   IP: 240.119.117.3
nginx-e2e-65999c5895-qd2zn   1/1     Running   fancni-e2e-node-1   IP: 240.119.94.2

Testing 240.119.117.4... 200   ← node-2 pod ✓
Testing 240.119.94.3...  200   ← node-1 pod ✓
Testing 240.119.117.3... 200   ← node-2 pod ✓
Testing 240.119.94.2...  200   ← node-1 pod ✓
```

✓ HTTP 200 from ALL pods, including cross-node.

---

## Summary

Cross-node connectivity is now working. Two issues were identified and fixed:

1. **Original fix (already in source):** Changed `-A FORWARD` (append) to `-I FORWARD 1` (insert at top) and added sysctls to disable ICMP redirects and rp_filter — this was correct but incomplete.

2. **Additional fix applied during this validation:** Changed `iptables` to `iptables-legacy` throughout the init script. On Ubuntu 24.04, `iptables` resolves to `iptables-nft`, but kube-proxy uses `iptables-legacy`. The two backends have separate rule tables, so rules inserted into `iptables-nft` were invisible to the packet forwarding path controlled by `iptables-legacy`.

## Files Changed

- `deploy/scripts/init-node.sh` — changed `iptables` → `iptables-legacy` for FORWARD and POSTROUTING rules
