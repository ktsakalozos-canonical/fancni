# Task: Quick validation of cross-node fix on existing debug VMs

## Context

Debug VMs `fancni-e2e-node-1` and `fancni-e2e-node-2` are still running from a previous e2e run (with `--no-cleanup`). They have a full Canonical K8s cluster with fancni installed and nginx-e2e deployment (4 pods). Cross-node pod connectivity was failing due to iptables rule ordering.

The fix has been applied to `deploy/scripts/init-node.sh`. We need to rebuild the init image, re-import it, and restart the DaemonSet pods to pick up the new init script.

## Objective

Validate that the iptables fix resolves cross-node pod connectivity without a full e2e re-run.

## Steps

1. **Rebuild only the init image** on the host:
   ```bash
   docker build -t fancni-init:latest -f deploy/docker/Dockerfile.init deploy/
   ```

2. **Re-import the init image** into both VMs:
   ```bash
   for VM in fancni-e2e-node-1 fancni-e2e-node-2; do
     docker save fancni-init:latest | lxc exec "$VM" -- /snap/k8s/current/bin/ctr --address /run/containerd/containerd.sock -n k8s.io images import -
   done
   ```

3. **Delete the fancni DaemonSet pods** to trigger recreation (init container will re-run with new script):
   ```bash
   lxc exec fancni-e2e-node-1 -- sudo k8s kubectl -n kube-system delete pods -l app=fancni
   ```

4. **Wait for DaemonSet pods to be Running again**:
   ```bash
   # Wait until all fancni pods show Running
   lxc exec fancni-e2e-node-1 -- sudo k8s kubectl -n kube-system get pods -l app=fancni -w
   ```

5. **Verify iptables rules are inserted at top** on both nodes:
   ```bash
   lxc exec fancni-e2e-node-1 -- iptables -L FORWARD -n --line-numbers | head -10
   lxc exec fancni-e2e-node-2 -- iptables -L FORWARD -n --line-numbers | head -10
   ```
   Expected: our ACCEPT rules for 240.x.x.0/24 appear at positions 1 and 2.

6. **Verify sysctls** on both nodes:
   ```bash
   lxc exec fancni-e2e-node-1 -- sysctl net.ipv4.conf.fan-240.send_redirects net.ipv4.conf.fan-240.rp_filter net.ipv4.conf.all.rp_filter
   lxc exec fancni-e2e-node-2 -- sysctl net.ipv4.conf.fan-240.send_redirects net.ipv4.conf.fan-240.rp_filter net.ipv4.conf.all.rp_filter
   ```
   Expected: all 0.

7. **Test cross-node connectivity**:
   ```bash
   # Get pod IPs
   lxc exec fancni-e2e-node-1 -- sudo k8s kubectl get pods -l app=nginx-e2e -o wide
   # Curl each pod from node-1
   POD_IPS=$(lxc exec fancni-e2e-node-1 -- sudo k8s kubectl get pods -l app=nginx-e2e -o jsonpath='{.items[*].status.podIP}')
   for IP in $POD_IPS; do
     echo "Testing $IP..."
     lxc exec fancni-e2e-node-1 -- curl -s -o /dev/null -w "%{http_code}" --max-time 5 "http://$IP"
     echo
   done
   ```
   Expected: HTTP 200 from ALL pods, including those on node-2.

## Non-goals
- Do not re-run the full e2e script
- Do not modify any source files
- Do not delete/recreate the nginx deployment

## Report
Report results (iptables output, sysctl values, curl results) in the completion report. Clearly state whether cross-node connectivity is now working.
