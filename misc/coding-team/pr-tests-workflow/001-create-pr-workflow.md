# Task: Create PR Tests Workflow

## Context
The project has no CI running on PRs. It's a Go 1.24 project with `make build`, `make test` (unit), and `make e2e` (LXC-based e2e). Existing workflows live in `.github/workflows/`.

## Objective
Create `.github/workflows/pr-tests.yml` that runs on PRs targeting `main` with two parallel jobs.

## Scope

Create a single file: `.github/workflows/pr-tests.yml`

### Job 1: `test`
- Runs on `ubuntu-latest`
- Steps: checkout, setup Go 1.24, `make build`, `make test`

### Job 2: `e2e`
- Runs on `ubuntu-latest`
- Steps:
  1. Checkout
  2. Setup Go 1.24 (needed for `make docker-build` which builds from source)
  3. Install LXD: `sudo snap install lxd` then `sudo lxd init --auto`
  4. Add the runner user to the `lxd` group and use `sg lxd` or `newgrp` to run commands in that group context — OR just use `sudo` for lxc commands. Since the e2e script uses `lxc` directly (no sudo), add the user to the lxd group and wrap the make e2e call with `sg lxd -c "make e2e"`.
  5. Install Helm: `sudo snap install helm --classic`
  6. Run `make e2e` (via `sg lxd` for group permissions)

### Trigger
```yaml
on:
  pull_request:
    branches: [main]
```

## Non-goals
- No caching, matrix builds, or artifact uploads
- No timeout tuning beyond a generous default (the e2e can take ~30min, so set timeout-minutes: 60 on that job)

## Constraints
- Nested virtualization is experimental on GitHub runners — this is accepted.
- The e2e script uses `lxc launch ... --vm` which requires KVM. If `/dev/kvm` doesn't exist the job will fail — that's acceptable for now.
