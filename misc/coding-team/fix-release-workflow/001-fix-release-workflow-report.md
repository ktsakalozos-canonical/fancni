# Completion Report: Fix release workflow networking and add manual trigger

## Summary
- Added `workflow_dispatch:` trigger to the `on:` block so the workflow can be run manually.
- Updated the job `if` condition to also pass on `workflow_dispatch` events (previously only `workflow_run` was handled).
- Added fallback `|| github.sha` to `ref:` (checkout) and `SHA` env var so manual dispatch uses the correct commit.
- Added a "Configure LXD network" step between `lxd waitready` and `rockcraft pack` to fix network access inside LXD containers on GitHub Actions runners.

## Files Changed
- `.github/workflows/release-latest.yml`

## Notable Tradeoffs / Risks
- None. Changes are minimal and additive; no existing steps were removed or renamed.
