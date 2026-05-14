# Task: Add retry-with-backoff to nightly workflow curl calls

## Context

Both `.github/workflows/nightly-dreams.yml` and `.github/workflows/nightly-knowledge.yml` call `https://models.github.ai/inference/chat/completions` via `curl -fsSL`. The `-f` flag causes curl to exit immediately on HTTP errors. The GitHub Models API returns 429 (rate limit) under load, which kills the workflow.

## Objective

Add a bash retry wrapper so that 429 responses are retried with exponential backoff, while other errors still fail fast.

## Scope

- `.github/workflows/nightly-dreams.yml` — 1 curl call (step "Call GitHub Models API (Morpheus)")
- `.github/workflows/nightly-knowledge.yml` — 5 curl calls (architect, developer, code-reviewer, code-reviewerer, repo-scout steps)

## Approach

In each workflow file, define a shell function early in the step (or in a preceding step that exports it) that wraps curl with retry logic. Something like:

```bash
retry_curl() {
  local max_retries=3
  local delay=10
  for i in $(seq 1 "$max_retries"); do
    HTTP_CODE=$(curl -sL -w '%{http_code}' -o /tmp/curl_response.json \
      "$@")
    if [ "$HTTP_CODE" -eq 200 ]; then
      cat /tmp/curl_response.json
      return 0
    elif [ "$HTTP_CODE" -eq 429 ]; then
      echo "::warning::Got 429 (attempt $i/$max_retries). Retrying in ${delay}s..."
      sleep "$delay"
      delay=$((delay * 2))
    else
      echo "::error::HTTP $HTTP_CODE from API"
      cat /tmp/curl_response.json >&2
      return 1
    fi
  done
  echo "::error::Exhausted retries (last HTTP $HTTP_CODE)"
  return 1
}
```

Then replace each `curl -fsSL ...` invocation with `retry_curl -H "Authorization: ..." -H "Content-Type: ..." -d "$PAYLOAD" <URL>`.

Key details:
- Drop the `-f` flag (the wrapper handles HTTP status itself via `-w '%{http_code}'`).
- Keep `-sL` (silent + follow redirects).
- The function outputs the response body to stdout on success, so `RESPONSE=$(retry_curl ...)` works as a drop-in replacement.
- For `nightly-knowledge.yml`, define the function once and reuse it across all 5 steps. Since each step runs in its own shell, the simplest approach is to either: (a) duplicate the function definition in each step, or (b) write it to a script file in an earlier step and source it. Option (a) is fine — use a YAML anchor/alias (`&retry_curl_fn` / `*retry_curl_fn`) or just paste it. Whichever is cleaner in practice.

## Non-goals

- Don't change API payloads, models, temperature, or max_tokens.
- Don't restructure workflow steps or parallelize calls.
- Don't add retries for non-curl failures (empty content, jq parsing, etc.) — those are already handled.

## Constraints

- Pure bash, no new actions or dependencies.
- Keep backoff reasonable (10s, 20s, 40s) — total worst-case wait ~70s per call site, acceptable for a nightly job.
