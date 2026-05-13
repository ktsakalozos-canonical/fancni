# Task: Fix review findings in nightly-dreams workflow

## Context
The workflow at `.github/workflows/nightly-dreams.yml` was reviewed. Two issues need fixing.

## Objective
Apply two targeted fixes to the workflow.

## Changes

### 1. Add API response validation (after the curl call)
After line 97 (the curl call), add a check:
- If `$RESPONSE` is empty, print an error with `::error::` and `exit 1`.
- Also check if `jq` can extract `.choices[0].message.content` — if the result is `null` or empty, error out.

### 2. Handle filename collision
In the "Write dream file" step, before writing, check if `$OUTPUT_FILE` already exists. If it does, append `-$(date +%H%M%S)` to `DREAM_WORD` to make it unique (e.g., `security-031500.md`).

## Non-goals
- Do not add retry logic, parameterization, or other changes beyond these two fixes.
