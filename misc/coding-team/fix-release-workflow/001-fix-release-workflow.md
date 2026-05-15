# Task: Fix release workflow networking and add manual trigger

## Context
The release workflow at `.github/workflows/release-latest.yml` fails because the LXD instance spawned by `rockcraft pack` has no network access inside GitHub Actions. Additionally, the workflow should be manually triggerable.

## Objective
Two changes to `.github/workflows/release-latest.yml`:

1. **Fix LXD networking**: Add a step *after* `lxd init --auto` / `lxd waitready` but *before* `rockcraft pack` that configures the network:
   ```yaml
   - name: Configure LXD network
     run: |
       sudo iptables -P FORWARD ACCEPT
       sudo lxc network set lxdbr0 ipv4.address=10.10.10.1/24 ipv4.nat=true
   ```

2. **Add manual trigger**: Add `workflow_dispatch:` to the `on:` block (no inputs needed).

3. **Fix the job `if` condition**: Currently it only passes for `workflow_run` events. Update it so the job also runs on `workflow_dispatch`. Simplest approach:
   ```yaml
   if: >-
     ${{ github.event_name == 'workflow_dispatch' ||
         (github.event.workflow_run.conclusion == 'success' &&
          github.event.workflow_run.event == 'push' &&
          github.event.workflow_run.head_branch == 'main') }}
   ```

4. **Fix checkout ref for manual dispatch**: The current `ref:` uses `github.event.workflow_run.head_sha` which is empty on manual dispatch. Use a fallback:
   ```yaml
   ref: ${{ github.event.workflow_run.head_sha || github.sha }}
   ```
   Similarly for the `SHA` env var:
   ```yaml
   SHA: ${{ github.event.workflow_run.head_sha || github.sha }}
   ```

## Non-goals
- Do not change rockcraft.yaml.
- Do not add workflow_dispatch inputs.
- Do not remove or rename existing steps.

## Constraints
- Keep the file concise; don't restructure beyond what's needed.
