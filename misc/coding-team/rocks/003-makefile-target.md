# Task: Add rock-build Makefile target

## Context
We added a `rockcraft.yaml` at the project root. Need a Makefile target to build it.

## Objective
Add a `rock-build` target to the Makefile that runs `rockcraft pack`.

## Scope
- File: `Makefile`
- Add `rock-build` to the `.PHONY` list
- Add the target after the existing `docker-build` target

## Constraints
- Do not modify existing targets.
