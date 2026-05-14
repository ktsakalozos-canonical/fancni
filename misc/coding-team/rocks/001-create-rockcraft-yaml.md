# Task: Create rockcraft.yaml

## Context
This project has two Docker images (CNI installer + node initializer). We're replacing them with a single rock that contains everything. The pod spec will select behavior via `command:` overrides.

## Objective
Create a `rockcraft.yaml` at the project root that produces one OCI rock image.

## Scope
- File: `/rockcraft.yaml` (project root)

## What the rock must contain
1. The `fancni` Go binary at `/fancni` (built from `./cmd/fancni/`)
2. `deploy/scripts/install-cni.sh` at `/install-cni.sh` (executable)
3. `deploy/scripts/init-node.sh` at `/init-node.sh` (executable)
4. APT packages: `iptables`, `iproute2`, `kmod`

## Key rockcraft.yaml details
- `name: fancni`
- `base: ubuntu@24.04`
- `build-base: ubuntu@24.04`
- `channel` for base/build-base snaps: use default (no override needed)
- `platforms: amd64`
- Use the `go` plugin for building the binary. Source is the project root; `go-builddir` is `cmd/fancni`. The built binary should end up at `/fancni` in the rock.
- Use a `dump` plugin part (or `overlay`) to stage the two shell scripts into the correct locations.
- Stage packages: `iptables`, `iproute2`, `kmod`
- No services defined (no default entrypoint — the orchestrator picks the script).

## Non-goals
- Do not modify existing Dockerfiles or scripts.
- Do not touch Makefile (separate task).

## Constraints
- Go version in `go.mod` is 1.24.13 — ensure the build snap or go plugin can satisfy this.
- Rock must be packable with `rockcraft pack` (though we won't run it in CI yet).
