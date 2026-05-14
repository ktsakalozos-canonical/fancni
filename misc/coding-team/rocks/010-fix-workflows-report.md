# Task 010: Fix Workflows — Report

## Summary

- **ci-tests.yml**: Renamed job display name from `E2e` to `E2E` (cosmetic fix).
- **release-latest.yml**: Added `Install LXD` and `Add runner to lxd group` steps before `Build rock`; wrapped `rockcraft pack` with `sg lxd -c "..."` so the runner has LXD group permissions.
- **release-latest.yml**: Removed the duplicate `env: REGISTRY:` block from the `Push rock to GHCR` step (it was already defined at the job level).

## Files Changed

- `.github/workflows/ci-tests.yml`
- `.github/workflows/release-latest.yml`

## Notable Tradeoffs / Risks

None. Changes are minimal and follow the exact instructions in the task brief.
