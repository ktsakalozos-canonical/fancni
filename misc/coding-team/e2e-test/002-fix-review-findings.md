# Task: Fix review findings in e2e test script

## Context

Code review found issues in `tests/e2e/test-e2e.sh`. This task fixes them.

## Fixes required

### 1. Critical: Wrong DaemonSet label selector (line 152)

The Helm chart DaemonSet uses label `app: fancni` (see `deploy/helm/fancni/templates/daemonset.yaml`).
The script uses `-l app.kubernetes.io/name=fancni` which will never match.

**Fix:** Change to `-l app=fancni`.

### 2. VM creation idempotency (lines 64-79)

The `--device root,size=20GiB` with fallback is fragile. Instead:
- Delete any pre-existing VMs at the start of Phase 1: `lxc delete --force "${NODE1}" "${NODE2}" 2>/dev/null || true`
- Use a single clean `lxc launch` without the device override fallback (LXC VMs get a 20GiB root by default anyway on most LXD setups; if you keep the device, use the correct syntax)

### 3. Bootstrap config: push file instead of heredoc-in-bash-c (lines 101-106)

The current `bash -c 'cat <<EOF ... EOF\nsudo k8s bootstrap...'` is fragile with shell escaping.

**Fix:** Write the config to a temp file on the host, push it into the VM, then run bootstrap:
```bash
local config="/tmp/fancni-bootstrap-config.yaml"
cat <<EOF > "${config}"
cluster-config:
  network:
    enabled: false
EOF
lxc file push "${config}" "${NODE1}/tmp/bootstrap-config.yaml"
rm -f "${config}"
lxc_exec "${NODE1}" sudo k8s bootstrap --file /tmp/bootstrap-config.yaml
```

### 4. Simplify k8s status wait (lines 108-109)

`k8s status --wait-ready` already blocks until ready or times out. Wrapping it in `wait_for` with polling is redundant.

**Fix:** Call it directly with a timeout:
```bash
lxc_exec "${NODE1}" sudo k8s status --wait-ready --timeout 300
```

## Non-goals
- Do NOT quote `${POD_IPS}` in the `for` loops — word splitting is intentional there
- Do NOT add anti-affinity to the nginx deployment
- Do NOT remove the `lxc_exec` helper
