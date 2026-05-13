# Report: Fix lxc file push needing pre-created directory

## Summary

- Added `lxc_exec "${NODE1}" mkdir -p /tmp/fancni-chart` before the `lxc file push` line in Phase 8 of `tests/e2e/test-e2e.sh`.
- This prevents the "Error: Not Found" failure when `/tmp/fancni-chart/` does not exist in the VM prior to the push.

## Files Changed

- `tests/e2e/test-e2e.sh`

## Notable Tradeoffs or Risks

None. The fix is minimal and safe; `mkdir -p` is idempotent.
