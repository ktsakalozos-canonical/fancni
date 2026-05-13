# Task 004: Add retry logic to snap install in e2e tests

## Context
The e2e test script (`tests/e2e/test-e2e.sh`) installs the `k8s` snap in Phase 3 (lines 88-89). Sometimes snapd is still restarting after cloud-init finishes, causing "unable to contact snap store" errors.

## Objective
Add retry logic around the `snap install` commands in Phase 3 so transient snapd/store connectivity issues are tolerated.

## Scope
- `tests/e2e/test-e2e.sh`, Phase 3 only (lines 87-89).

## Design
Add a helper function (or inline loop) that retries the snap install up to 5 times with a 30-second delay between attempts. Apply it to both `snap install` calls. Log each retry attempt.

Example approach — add a `retry` helper near the existing `wait_for` helper:

```bash
# retry <max_attempts> <delay_seconds> <cmd…>
retry() {
  local max="$1"; shift
  local delay="$1"; shift
  local attempt=1
  while true; do
    if "$@"; then
      return 0
    fi
    if (( attempt >= max )); then
      log "Command failed after ${max} attempts: $*"
      return 1
    fi
    log "  Attempt ${attempt}/${max} failed, retrying in ${delay}s..."
    sleep "${delay}"
    (( attempt++ ))
  done
}
```

Then Phase 3 becomes:
```bash
retry 5 30 lxc_exec "${NODE1}" snap install k8s --classic --channel=1.35-classic/stable
retry 5 30 lxc_exec "${NODE2}" snap install k8s --classic --channel=1.35-classic/stable
```

## Non-goals
- Do not change any other phases.
- Do not modify the wait_for helper.
