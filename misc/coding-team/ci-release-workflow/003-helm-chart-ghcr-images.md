# Task 003: Helm chart default images point to GHCR

## Context
The release workflow now pushes images to `ghcr.io/ktsakalozos-canonical/fancni` and `ghcr.io/ktsakalozos-canonical/fancni-init`. The helm chart's `values.yaml` currently uses bare names (`fancni`, `fancni-init`) as repository values.

## Objective
Update `deploy/helm/fancni/values.yaml` so the default image repositories include the full GHCR path.

## Changes
In `deploy/helm/fancni/values.yaml`:
- `images.cni.repository`: `fancni` → `ghcr.io/ktsakalozos-canonical/fancni`
- `images.init.repository`: `fancni-init` → `ghcr.io/ktsakalozos-canonical/fancni-init`

Tag stays `latest`. pullPolicy stays `IfNotPresent`.

## Non-goals
- Do not change templates or any other values.
