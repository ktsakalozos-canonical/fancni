# Report: 002-harden-retry-edge-cases

## Summary

- Added a `HTTP_CODE` empty/`000` guard immediately after each `curl` call in all 6 `retry_curl` instances. Connection-level failures (DNS, reset, etc.) are now treated as retryable with backoff, matching the 429 behaviour.
- Replaced `/tmp/curl_response.json` with `/tmp/curl_response_$$.json` (PID-scoped) in all 6 instances across both workflow files to eliminate theoretical temp-file collisions between concurrent runs.

## Files Changed

- `.github/workflows/nightly-dreams.yml` — 1 instance updated
- `.github/workflows/nightly-knowledge.yml` — 5 instances updated

## Notable Tradeoffs / Risks

None. Changes are purely defensive and backward-compatible; the normal success/429/error paths are unchanged.
