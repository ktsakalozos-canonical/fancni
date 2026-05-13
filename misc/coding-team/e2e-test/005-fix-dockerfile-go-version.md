# Task: Fix Dockerfile.cni Go version mismatch

## Context

`deploy/docker/Dockerfile.cni` uses `golang:1.22-bookworm` but `go.mod` requires `go 1.24.13`. The Docker build fails with "go.mod requires go >= 1.24.13 (running go 1.22.12)".

## Fix

In `deploy/docker/Dockerfile.cni`, change line 2:
```
FROM golang:1.22-bookworm AS builder
```
to:
```
FROM golang:1.24-bookworm AS builder
```
