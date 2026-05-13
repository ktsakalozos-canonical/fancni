# Task 005 Completion Report: Init Container & Dockerfiles

## Summary

- Created `deploy/docker/Dockerfile.cni`: multi-stage build using `golang:1.22-bookworm` to compile the CNI binary, then copies it into a minimal `busybox:1.36` install image along with `install-cni.sh`.
- Created `deploy/docker/Dockerfile.init`: `ubuntu:24.04` image with `iptables`, `iproute2`, and `kmod` installed; runs `init-node.sh` as entrypoint.
- Created `deploy/scripts/init-node.sh` and `deploy/scripts/install-cni.sh` matching the exact logic specified in the task brief; both use `set -e` and are idempotent.
- Added `docker-build-cni`, `docker-build-init`, and `docker-build` targets to the Makefile.

## Files Changed

- `deploy/docker/Dockerfile.cni` (new)
- `deploy/docker/Dockerfile.init` (new)
- `deploy/scripts/init-node.sh` (new)
- `deploy/scripts/install-cni.sh` (new)
- `Makefile` (updated `.PHONY` and added three docker build targets)

## Notable Tradeoffs or Risks

- Both Dockerfiles use the repo root as build context (required for the multi-stage copy of scripts); this is consistent with the Makefile targets using `.` as the context.
- `install-cni.sh` uses `exec sleep infinity` at the end to keep the DaemonSet pod alive, as specified; this is intentional per the brief.
- `go build ./...` and `go test ./...` both pass unchanged.
