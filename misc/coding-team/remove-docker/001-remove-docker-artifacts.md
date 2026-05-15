# Task: Remove Docker artifacts

## Context
The project has migrated to rocks (rockcraft). The old Docker build path is dead code that causes confusion.

## Objective
Remove all Docker build artifacts and update docs to reflect the rock-based workflow.

## Scope

1. **Delete** `deploy/docker/Dockerfile.cni` and `deploy/docker/Dockerfile.init`, then remove the now-empty `deploy/docker/` directory.

2. **Makefile** — remove the three docker targets (`docker-build-cni`, `docker-build-init`, `docker-build`) and their `.PHONY` entries. Leave everything else intact.

3. **`docs/development.md`** — replace the `make docker-build` reference with `make rock-build` and update the comment to say "Build rock image" (or similar).

## Non-goals
- Do not touch anything under `misc/`.
- Do not change the e2e test, CI workflows, or release workflow.
- Do not modify `rockcraft.yaml` or the rock-build target.
