# Completion Report: Remove Docker Artifacts

## Summary

- Deleted `deploy/docker/Dockerfile.cni`, `deploy/docker/Dockerfile.init`, and the now-empty `deploy/docker/` directory.
- Removed `docker-build-cni`, `docker-build-init`, and `docker-build` targets (and their `.PHONY` entries) from the root `Makefile`.
- Updated `docs/development.md` to replace `make docker-build` with `make rock-build` and updated the comment to "Build rock image".

## Files Changed

- `deploy/docker/Dockerfile.cni` — deleted
- `deploy/docker/Dockerfile.init` — deleted
- `deploy/docker/` — directory removed
- `Makefile` — docker targets and `.PHONY` entries removed
- `docs/development.md` — docker-build reference replaced with rock-build

## Notable Tradeoffs or Risks

None. Changes are purely deletions/cleanup; no logic was altered.
