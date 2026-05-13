# Task 002: Release Latest Workflow — Completion Report

## Summary

- Created `.github/workflows/release-latest.yml` that triggers via `workflow_run` on `ci-tests` completion.
- The job only runs when the triggering workflow succeeded, was a push event, and targeted the `main` branch.
- Builds both `Dockerfile.cni` (fancni) and `Dockerfile.init` (fancni-init) images and pushes each with two tags: commit SHA and `latest`, to GHCR.
- Uses `actions/checkout@v4` with `ref: head_sha` and `docker/login-action@v3` as specified.

## Files Changed

- `.github/workflows/release-latest.yml` (new)

## Notable Tradeoffs or Risks

- No multi-arch builds or layer caching, per the brief's non-goals.
- The `if` condition uses the `${{ }}` expression syntax inline with the YAML block scalar; this is intentional and matches GitHub Actions semantics for `workflow_run` triggers.
