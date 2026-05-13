# Task 005 – Fix Dockerfile Go Version: Completion Report

## Summary
- Updated `deploy/docker/Dockerfile.cni` builder stage from `golang:1.22-bookworm` to `golang:1.24-bookworm` to align the Docker build environment with the Go version used in CI and the rest of the project.

## Files Changed
- `deploy/docker/Dockerfile.cni`

## Notable Tradeoffs or Risks
- None. This is a single-line bump to a newer patch-compatible Go release image; no code changes were required.
