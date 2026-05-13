# Task Report: Include Prior Knowledge in Nightly Generation Prompts

## Summary

- For each of the 5 "Generate X knowledge" steps in `.github/workflows/nightly-knowledge.yml`, added a `PRIOR_KNOWLEDGE` read before the API call using `cat misc/knowledge/<agent>-knowledge.md 2>/dev/null || echo ""` — handles missing file gracefully on first run.
- When prior knowledge exists, it is appended to the user message under an `=== Prior Knowledge ===` section; otherwise the user message is unchanged.
- Added a line to each system prompt: "You may refine, update, or build on the prior knowledge provided — keep what is still accurate and update what has changed."
- The `--arg user` in the `jq` PAYLOAD call was updated from `$PROJECT_CONTEXT` to `$USER_MESSAGE` in all 5 steps.

## Files Changed

- `.github/workflows/nightly-knowledge.yml`

## Notable Tradeoffs or Risks

- None. Changes are minimal and additive; no restructuring of the workflow.
