# Task: Harden retry_curl edge cases

## Context

The `retry_curl` function was added in the previous task. Two edge cases were flagged in review.

## Changes needed

In all 6 instances of `retry_curl` across both workflow files:

1. **Handle non-HTTP curl failures.** If curl itself fails (DNS, connection reset, etc.), `HTTP_CODE` may be empty or `000`. Add a check right after the curl call:
   ```bash
   if [ -z "$HTTP_CODE" ] || [ "$HTTP_CODE" = "000" ]; then
     echo "::warning::curl failed (no HTTP response, attempt $i/$max_retries). Retrying in ${delay}s..."
     sleep "$delay"
     delay=$((delay * 2))
     continue
   fi
   ```
   This treats connection-level failures as retryable (same as 429).

2. **Use PID-scoped temp file.** Replace `/tmp/curl_response.json` with `/tmp/curl_response_$$.json` to avoid any theoretical collision.

## Non-goals

- Don't change anything else about the function or workflow logic.
