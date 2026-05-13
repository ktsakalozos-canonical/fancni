# Report: 006-fix-containerd-socket

## Summary

- Changed `CTR_SOCK` in `tests/e2e/test-e2e.sh` from `/var/snap/k8s/common/run/containerd.sock` to `/run/containerd/containerd.sock`.
- The old path was snap-specific; the new path is the standard containerd socket location used by the test environment.

## Files Changed

- `tests/e2e/test-e2e.sh`

## Notable Tradeoffs or Risks

None. Single-line config change with no logic impact.
