# Task: Fix timeout duration format

## Context

`k8s status --wait-ready --timeout 300` fails with: `invalid argument "300" for "--timeout" flag: time: missing unit in duration "300"`. The flag expects a Go duration string.

## Fix

In `tests/e2e/test-e2e.sh`, line 105, change:
```
lxc_exec "${NODE1}" sudo k8s status --wait-ready --timeout 300
```
to:
```
lxc_exec "${NODE1}" sudo k8s status --wait-ready --timeout 300s
```

That's the only change needed.
