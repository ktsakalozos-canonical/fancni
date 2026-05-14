# Report: 009-update-release

## Summary
- Removed `docker/login-action` step and both Docker build+push steps (fancni and fancni-init)
- Added `Install Rockcraft`, `Build rock`, and `Push rock to GHCR` steps using `rockcraft pack` and `rockcraft.skopeo`
- Added `name: Release` to the job for capitalized display name
- Auth is now handled inline via `--dest-creds` in skopeo; no Docker login action needed

## Files changed
- `.github/workflows/release-latest.yml`

## Notable tradeoffs or risks
- The rock filename `fancni_0.1_amd64.rock` is hardcoded; if the version in `rockcraft.yaml` changes, this will need updating.
- `rockcraft.skopeo` is bundled with the rockcraft snap, so no additional install step is required.
