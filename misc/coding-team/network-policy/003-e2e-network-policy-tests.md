# Task: E2E network policy tests + flip default + update plan.md

## Context

fancni now supports network policies via kube-router (firewall-only mode) deployed as a second DaemonSet. The Helm chart has `networkPolicy.enabled: false` as default. We want to:
1. Change the default to `true` (network policies on by default)
2. Extend the e2e test to verify network policy enforcement
3. Update `misc/plan.md` to reflect that network policies are now in scope

## Objective

### 1. Flip default in values.yaml

In `deploy/helm/fancni/values.yaml`, change:
```yaml
networkPolicy:
  enabled: false
```
to:
```yaml
networkPolicy:
  enabled: true
```

### 2. Extend e2e test (`tests/e2e/test-e2e.sh`)

Add the following after the existing Phase 12 (before the summary), keeping the same style/conventions:

**Phase 12: Transfer kube-router image into VMs**

The kube-router image (`cloudnativelabs/kube-router:v2.9.0`) needs to be pulled on the host and transferred into both VMs, same pattern as nginx/fancni images in Phase 7. Add it to the existing image transfer loop in Phase 7 (just add `"cloudnativelabs/kube-router:v2.9.0"` to the list and add a `docker pull` for it alongside nginx).

**Phase 12 (renumber existing summary to Phase 15): Verify kube-router pods running**

Wait for pods with label `app=kube-router` in `kube-system` to be Running on both nodes.

**Phase 13: Test deny-all NetworkPolicy blocks traffic**

Steps:
1. Create a namespace `netpol-test`
2. Deploy a simple nginx pod (1 replica) in `netpol-test` with label `app=web`, wait for Running
3. Deploy a client pod (using `nginx:latest` image, just need curl) in `netpol-test` with label `app=client`, wait for Running
4. Verify connectivity works BEFORE policy: `kubectl exec client-pod -- curl --max-time 5 http://<web-pod-ip>` should return 200
5. Apply a deny-all ingress NetworkPolicy targeting `app=web`:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: netpol-test
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
```
6. Wait 10 seconds for kube-router to sync rules
7. Verify connectivity is NOW BLOCKED: `kubectl exec client-pod -- curl --max-time 5 http://<web-pod-ip>` should fail (non-200 or timeout)

**Phase 14: Test allow NetworkPolicy restores traffic**

Steps:
1. Apply an allow policy permitting ingress from `app=client`:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-client
  namespace: netpol-test
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: client
```
2. Wait 10 seconds for sync
3. Verify connectivity is restored: curl from client pod to web pod returns 200

**Phase 15: Summary (renumbered from 12)**

Add network policy test results to the summary output. Fail the overall test if any network policy assertion fails.

### 3. Update `misc/plan.md`

Make these changes:
- In "Constraints & Assumptions" section, remove the line `- **No network policies**: Out of scope`
- In "Out of Scope" section, remove the line `- Network policy enforcement`
- Add a new section **after** "Implementation Phases" (before "Constraints & Assumptions") titled `## Network Policy Support` with this content:

```markdown
## Network Policy Support

Network policies are enforced via [kube-router](https://github.com/cloudnativelabs/kube-router) running in firewall-only mode (`--run-firewall=true`). kube-router is deployed as a second DaemonSet within the same Helm chart.

- **Enabled by default** (`networkPolicy.enabled: true` in values.yaml)
- Watches `NetworkPolicy` resources via the Kubernetes API
- Enforces ingress/egress rules using iptables + ipsets on each node
- When enabled, fancni's init script skips blanket FORWARD ACCEPT rules (kube-router manages the FORWARD chain)
- The MASQUERADE rule for pod egress remains active regardless

Configuration in `values.yaml`:
```yaml
networkPolicy:
  enabled: true
  image:
    repository: cloudnativelabs/kube-router
    tag: v2.9.0
    pullPolicy: IfNotPresent
  iptablesSyncPeriod: "5m"
```
```

## Constraints

- Keep the same e2e test style: `log()`, `wait_for()`, `lxc_exec()` helpers
- Use `sudo k8s kubectl` for all kubectl commands (that's how the test accesses the cluster)
- All kubectl apply should use heredoc piped to `lxc_exec NODE1 sudo k8s kubectl apply -f -`
- Renumber existing Phase 12 (Summary) appropriately
- The kube-router sync period is 5m by default but the test should work because we wait for connectivity to fail/succeed with retries

## Non-goals

- Do not modify the kube-router DaemonSet template
- Do not add unit tests
- Do not modify the Makefile
