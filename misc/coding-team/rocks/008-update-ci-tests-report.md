# Report: Update ci-tests.yml for rockcraft

## Summary
- Added `name: Test` and `name: E2e` display name fields to the `test` and `e2e` jobs respectively.
- Removed the "Setup Go" step from the `e2e` job (Go is now handled internally by rockcraft).
- Added "Install Rockcraft" step (`sudo snap install rockcraft --classic`) in the `e2e` job, immediately after the Checkout step.

## Files Changed
- `.github/workflows/ci-tests.yml`

## Notable Tradeoffs or Risks
- None. The `test` job retains its "Setup Go" step unchanged, as it still runs unit tests directly with Go.
