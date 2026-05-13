# Task: Fix nginx image pull timeout in e2e test

## Context

Phase 9 times out (120s) waiting for nginx pods because `nginx:latest` must be pulled from Docker Hub inside LXC VMs on a cold cache. fancni images are pre-loaded but nginx is not.

## Objective

Make nginx image available without relying on in-VM Docker Hub pulls.

## Changes to `tests/e2e/test-e2e.sh`

1. **Phase 7**: Add `nginx:latest` to the image transfer loop. Pull it on the host first if needed:
   ```bash
   docker pull nginx:latest
   ```
   Then add `"nginx:latest"` to the `for IMAGE in ...` loop.

2. **Phase 9**: Increase wait timeout from `120` to `300` (safety net for any remaining slow operations).

## Non-goals
- Don't restructure the test phases
- Don't add diagnostic dumps (keep it simple for now)
