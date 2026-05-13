# Task: Fix lxc file push needing pre-created directory

## Context

`lxc file push --recursive ... ${NODE1}/tmp/fancni-chart/` fails with "Error: Not Found" when `/tmp/fancni-chart/` doesn't exist in the VM.

## Fix

In `tests/e2e/test-e2e.sh`, before the `lxc file push` line in Phase 8, add a `mkdir -p`:

Change:
```bash
lxc file push --recursive "${REPO_ROOT}/deploy/helm/fancni/" "${NODE1}/tmp/fancni-chart/"
```
to:
```bash
lxc_exec "${NODE1}" mkdir -p /tmp/fancni-chart
lxc file push --recursive "${REPO_ROOT}/deploy/helm/fancni/" "${NODE1}/tmp/fancni-chart/"
```
