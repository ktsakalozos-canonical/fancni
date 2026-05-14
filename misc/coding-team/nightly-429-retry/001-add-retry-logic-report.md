# Completion Report: Add retry logic to nightly workflow curl calls

## Summary

- Replaced all `curl -fsSL` invocations in both nightly workflow files with a `retry_curl` bash function that retries up to 3 times on HTTP 429, using exponential backoff (10s → 20s → 40s). All other non-200 HTTP codes still fail immediately.
- The `-f` flag was dropped; HTTP status is now checked explicitly via `-w '%{http_code}'`. `-sL` (silent + follow redirects) is retained.
- For `nightly-knowledge.yml` (5 API call steps), the function definition is duplicated in each step, since each step runs in its own shell.
- For `nightly-dreams.yml` (1 API call step), the function is defined inline in the same step before the call.

## Files Changed

- `.github/workflows/nightly-dreams.yml`
- `.github/workflows/nightly-knowledge.yml`

## Notable Tradeoffs

- The function is duplicated 5 times in `nightly-knowledge.yml` rather than written to a script file and sourced. This is slightly verbose but avoids adding a new file to the repo and keeps each step self-contained, per the task brief's suggestion.
