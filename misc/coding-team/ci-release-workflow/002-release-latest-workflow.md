# Task 002: Release latest workflow

## Context
CI workflow is now called `ci-tests` and runs on push to main. We need a release workflow that triggers after CI succeeds on main, builds both container images, and pushes them to GHCR.

## Objective
Create `.github/workflows/release-latest.yml` that:
1. Triggers via `workflow_run` on `ci-tests` completion.
2. Only runs if the triggering workflow succeeded, was a push event, and targeted main.
3. Builds both Docker images (`Dockerfile.cni` and `Dockerfile.init`) and pushes to GHCR with two tags each: the commit SHA and `latest`.

## Scope
- New file: `.github/workflows/release-latest.yml`

## Design details

```yaml
name: release-latest

on:
  workflow_run:
    workflows: ["ci-tests"]
    types: [completed]

jobs:
  release:
    if: >-
      ${{ github.event.workflow_run.conclusion == 'success' &&
          github.event.workflow_run.event == 'push' &&
          github.event.workflow_run.head_branch == 'main' }}
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    env:
      REGISTRY: ghcr.io/${{ github.repository_owner }}
      SHA: ${{ github.event.workflow_run.head_sha }}
    steps:
      - Checkout (at the SHA from workflow_run)
      - Log in to GHCR (docker/login-action, registry ghcr.io, username ${{ github.actor }}, password ${{ secrets.GITHUB_TOKEN }})
      - Build and push fancni image:
          docker build -t $REGISTRY/fancni:$SHA -t $REGISTRY/fancni:latest -f deploy/docker/Dockerfile.cni .
          docker push $REGISTRY/fancni:$SHA
          docker push $REGISTRY/fancni:latest
      - Build and push fancni-init image:
          docker build -t $REGISTRY/fancni-init:$SHA -t $REGISTRY/fancni-init:latest -f deploy/docker/Dockerfile.init .
          docker push $REGISTRY/fancni-init:$SHA
          docker push $REGISTRY/fancni-init:latest
```

Use `actions/checkout@v4` with `ref: ${{ github.event.workflow_run.head_sha }}` to ensure we build the correct commit.

## Non-goals
- No complex verification script (the `if` condition is sufficient since we only have one CI workflow).
- No multi-arch builds.
- No caching.

## Constraints
- Workflow name in the file must be `release-latest`.
- Use `docker/login-action@v3` for GHCR login.
