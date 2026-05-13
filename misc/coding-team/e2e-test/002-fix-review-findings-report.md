# Completion Report: Fix review findings in e2e test script

## Summary

- **Fix 1 (label selector):** Changed `-l app.kubernetes.io/name=fancni` to `-l app=fancni` on line 152 (now ~line 148) to match the actual Helm chart DaemonSet label.
- **Fix 2 (VM idempotency):** Added `lxc delete --force "${NODE1}" "${NODE2}" 2>/dev/null || true` before VM creation and removed the fragile `--device root,size=20GiB` fallback pattern; each VM is now launched with a single clean `lxc launch` command.
- **Fix 3 (bootstrap config):** Replaced `bash -c 'cat <<EOF ... EOF\nsudo k8s bootstrap...'` with a host-side temp file written via heredoc, pushed into the VM with `lxc file push`, then deleted; bootstrap runs via `lxc_exec`.
- **Fix 4 (k8s status wait):** Removed the redundant `wait_for` wrapper around `k8s status --wait-ready`; now calls `lxc_exec "${NODE1}" sudo k8s status --wait-ready --timeout 300` directly. Note: `local config=` was corrected to `config=` since the code is at script top-level (not inside a function).

## Files changed

- `tests/e2e/test-e2e.sh`

## Notable tradeoffs or risks

- Removing `--device root,size=20GiB` relies on the LXD default root disk size being sufficient (typically 10–30 GiB depending on the profile). If the environment uses a very small default pool, this could be an issue, but the task brief explicitly requested this simplification.
