# Task 001: CI workflow that also runs on push to main

## Context
Currently `pr-tests.yml` only triggers on `pull_request`. We need a CI workflow that also runs on `push` to `main` so the release workflow can use `workflow_run` to gate on it.

## Objective
Rename `pr-tests.yml` to `ci-tests.yml` and add `push: branches: [main]` to its trigger, keeping the existing `pull_request` trigger. Update the workflow `name:` to `ci-tests` (this name is referenced by the release workflow later).

## Scope
- `.github/workflows/pr-tests.yml` → rename to `.github/workflows/ci-tests.yml`
- Update `name:` field to `ci-tests`
- Add push trigger for main branch

## Non-goals
- Do not change the job steps, Go version, or test commands.
- Do not create the release workflow yet.

## Constraints
- The workflow name must be exactly `ci-tests` (will be referenced by workflow_run in a later task).
