# Report: Fix HTTP assertion with retries and output bug

## Summary

- Replaced the single-attempt curl call in Phase 11 with a retry loop (up to 5 attempts, 10s sleep between each) to tolerate fan tunnel convergence delay after DaemonSet startup.
- Fixed the doubled output bug (`000000`) by moving `|| true` outside the command substitution and initializing `HTTP_CODE="000"` before the loop instead of using `|| echo "000"` inside `$()`.
- Pass/fail logic moved outside the inner retry loop so it runs once after all attempts are exhausted.

## Files Changed

- `tests/e2e/test-e2e.sh` — Phase 11 HTTP assertion loop rewritten.

## Tradeoffs / Risks

- Worst-case added latency per pod is 4 × 10s = 40s if all retries fail. With 4 pods that's up to ~160s extra on failure paths. Acceptable for an E2E test.
