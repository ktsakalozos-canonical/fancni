# Task 004 Completion Report: Add retry logic to snap install in e2e tests

## Summary
- Added a `retry <max_attempts> <delay_seconds> <cmd…>` helper function in the Helpers section of `test-e2e.sh`, placed after `wait_for` and before `lxc_exec`.
- Updated Phase 3 to wrap both `snap install` calls with `retry 5 30`, tolerating up to 5 attempts with 30-second delays between retries.
- Each failed attempt is logged; if all attempts are exhausted, the command is reported as failed and the script exits non-zero (via `set -e`).

## Files Changed
- `tests/e2e/test-e2e.sh`

## Notable Tradeoffs or Risks
- None. Change is minimal and scoped to Phase 3 only, as required.
