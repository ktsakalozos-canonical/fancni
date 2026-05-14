# Task 007: Fix Rock Import — Completion Report

## Summary
- Added `--base-name "ghcr.io/ktsakalozos-canonical/fancni"` to the `ctr images import` call in Phase 7 of the e2e test, so the imported image gets a deterministic name instead of a timestamp-based one.
- Updated the subsequent `images tag` source reference from `fancni:0.1` to `ghcr.io/ktsakalozos-canonical/fancni:0.1` to match the new import name.

## Files Changed
- `tests/e2e/test-e2e.sh`

## Notable Tradeoffs / Risks
None. Change is minimal and exactly as specified in the task brief.
