# Task: Update ci-tests.yml for rockcraft

## Context
The e2e test now uses `make rock-build` (rockcraft) instead of `make docker-build`. The CI workflow needs rockcraft installed in the e2e job. Go setup is no longer needed in the e2e job since rockcraft builds Go internally.

## Objective
Update `.github/workflows/ci-tests.yml` to support rockcraft-based e2e and fix job naming.

## Changes

### Job names — capitalize
- `test` job: rename to `Test`
- `e2e` job: rename to `E2e`

These are the display names (the `name:` field under `jobs.<id>`). Keep the job IDs lowercase (the YAML keys `test:` and `e2e:` stay as-is for cross-workflow references).

Wait — GitHub Actions jobs don't have a separate `name:` by default. The display name IS the key. To capitalize the display name, add a `name:` field:
- Under `jobs.test:` add `name: Test`
- Under `jobs.e2e:` add `name: E2e`

### e2e job changes
1. **Remove** the "Setup Go" step (the rockcraft snap handles Go build internally).
2. **Add** a step to install rockcraft after the Checkout step:
   ```yaml
   - name: Install Rockcraft
     run: sudo snap install rockcraft --classic
   ```

## Non-goals
- Do not change the `test` job steps (it still uses Go directly for unit tests).
- Do not change LXD setup, Docker/LXD nftables fix, or any other existing steps.
