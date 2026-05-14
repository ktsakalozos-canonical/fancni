# Task: Fix workflow issues

## Context
Review identified that the release job needs LXD for `rockcraft pack`, and minor cosmetic issues.

## Changes

### ci-tests.yml
- Change `name: E2e` to `name: E2E`

### release-latest.yml
1. Add LXD install + init steps before "Build rock":
   ```yaml
   - name: Install LXD
     run: |
       sudo snap install lxd
       sudo lxd init --auto
       sudo lxd waitready

   - name: Add runner to lxd group
     run: sudo usermod -aG lxd "$USER"
   ```

2. Wrap `rockcraft pack` with `sg lxd` so the runner has LXD group permissions:
   ```yaml
   - name: Build rock
     run: sg lxd -c "rockcraft pack"
   ```

3. Remove the duplicate `env: REGISTRY:` block from the "Push rock to GHCR" step (it's already defined at the job level).

## Non-goals
- Do not change anything else.
