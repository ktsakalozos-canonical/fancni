# Completion Report: 001-create-pr-workflow

## Summary

- Created `.github/workflows/pr-tests.yml` with two parallel jobs (`test` and `e2e`) triggered on PRs targeting `main`.
- `test` job runs on `ubuntu-latest`: checkout → setup Go 1.24 → `make build` → `make test`.
- `e2e` job runs on `ubuntu-latest` with `timeout-minutes: 60`: checkout → setup Go 1.24 → install LXD + init → add runner to `lxd` group → install Helm → `sg lxd -c "make e2e"` for correct group context.

## Files Changed

- `.github/workflows/pr-tests.yml` (new file)

## Notable Tradeoffs / Risks

- LXC VM-based e2e requires KVM (`/dev/kvm`). On standard GitHub-hosted runners this may not be available; the task brief explicitly accepts this as a known limitation.
- No caching or artifact uploads per the brief's non-goals.
