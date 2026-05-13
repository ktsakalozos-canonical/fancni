# Task 005 – Completion Report

## Summary
- Added 9 missing Helm values to `docs/configuration.md` to match `deploy/helm/fancni/values.yaml`
- Values were grouped logically: `images.*` pullPolicy entries placed adjacent to their respective repo/tag entries; `cni.version` placed after `cni.confFileName`; `tolerations` and `networkPolicy.*` appended as new groups at the end
- Existing entries and descriptions were not modified

## Files Changed
- `docs/configuration.md`

## Notable Tradeoffs / Risks
- None. Change is purely additive documentation.
