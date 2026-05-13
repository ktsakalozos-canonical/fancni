# Task: Fix cross-node pod connectivity (iptables rule ordering + sysctls)

## Context

Cross-node pod traffic fails in the e2e test. Root cause: `init-node.sh` appends (`-A`) iptables FORWARD rules, but Canonical K8s inserts a `KUBE-FORWARD` chain jump early in FORWARD that contains `DROP ctstate INVALID`. Fan overlay return traffic hits this DROP before our ACCEPT rules are ever evaluated.

Secondary issue: the kernel sends ICMP redirects on the fan bridge (because packets arrive and leave on the same interface), which is noisy and could confuse some stacks.

## Objective

Make cross-node fan overlay traffic survive the KUBE-FORWARD chain by ensuring our ACCEPT rules are evaluated first.

## Changes required in `deploy/scripts/init-node.sh`

1. **Change `-A FORWARD` to `-I FORWARD 1`** for both ACCEPT rules (lines 42-44).
   - The `-C` check still works (checks existence regardless of position).
   - Use `-I FORWARD 1` to insert at the top of the chain.

2. **Add sysctl settings** after the bridge is created (after the `fanctl up` block, around line 34):
   ```bash
   # Disable ICMP redirects on fan bridge (packets legitimately arrive and leave on same interface)
   nsenter --target 1 --net -- sysctl -w net.ipv4.conf.${BRIDGE_NAME}.send_redirects=0
   # Disable reverse path filtering on fan bridge and tunnel
   nsenter --target 1 --net -- sysctl -w net.ipv4.conf.${BRIDGE_NAME}.rp_filter=0
   nsenter --target 1 --net -- sysctl -w net.ipv4.conf.all.rp_filter=0
   ```

3. **Enable IP forwarding explicitly** (belt-and-suspenders, k8s should set this but be safe):
   ```bash
   nsenter --target 1 --net -- sysctl -w net.ipv4.ip_forward=1
   ```

## Non-goals
- Do not change the MASQUERADE rule (it's fine as `-A` in nat/POSTROUTING)
- Do not fix the duplicate MAC issue (separate task)
- Do not modify the e2e test script

## Constraints
- Keep the `-C` idempotency check. It checks rule existence regardless of position, so this is fine.
- The `nsenter --target 1 --net` prefix is required for all network operations (we're in a container).
