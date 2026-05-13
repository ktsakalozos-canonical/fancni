# Task: Fix HTTP assertion with retries and output bug

## Context

Cross-node HTTP checks fail because fan networking needs time to converge after the DaemonSet starts. Also the curl error handling produces "000000" (double output) instead of "000".

## Fixes

### 1. Add retries to HTTP checks

Wrap each pod's HTTP check in a retry loop (e.g., up to 5 attempts with 10s sleep between). The fan tunnel between nodes may not be immediately ready.

Replace the Phase 11 loop with something like:

```bash
for IP in ${POD_IPS}; do
  HTTP_CODE="000"
  for attempt in $(seq 1 5); do
    HTTP_CODE=$(lxc_exec "${NODE1}" \
      curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${IP}" 2>/dev/null) || true
    if [[ "${HTTP_CODE}" == "200" ]]; then
      break
    fi
    log "  Attempt ${attempt}/5 for ${IP} returned ${HTTP_CODE}, retrying in 10s..."
    sleep 10
  done
  # ... rest of pass/fail logic
done
```

Key changes:
- The `|| true` goes AFTER the command substitution closing paren to avoid appending to the captured output
- Retry up to 5 times with 10s between attempts
- Log retry attempts

### 2. Fix curl output capture

The current code:
```bash
HTTP_CODE=$(lxc_exec "${NODE1}" \
    curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${IP}" 2>/dev/null || echo "000")
```

The `|| echo "000"` is inside the `$()`, so if curl outputs "000" (connection refused) AND then the `echo "000"` runs, you get "000000".

Fix: use `|| true` instead, and initialize HTTP_CODE to "000" before the subshell. If curl fails completely, HTTP_CODE stays "000" from the `|| true` path. Or simply: capture the output with `|| true` and if the result is empty, set it to "000".
