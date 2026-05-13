# Completion Report: 001-nightly-dreams-workflow

## Summary

- Created `.github/workflows/nightly-dreams.yml` — a nightly cron workflow (`0 3 * * *`) that checks out the repo, gathers project context (tree, recent commits, `go.mod`, `Makefile`), calls the GitHub Models API as "Morpheus" (`openai/gpt-4.1`), writes the response to `misc/nightly-dreams/<word>.md`, and commits it back with `[skip ci]`.
- Created `misc/nightly-dreams/.gitkeep` so the output directory is tracked in git before the first dream is written.
- The workflow exits cleanly (exit 0) when there are no commits in the last 7 days, per the brief's graceful-skip requirement.
- Only `actions/checkout` is used as an external action; everything else is shell, `curl`, and `jq`.

## Files Changed

- `.github/workflows/nightly-dreams.yml` (new)
- `misc/nightly-dreams/.gitkeep` (new)

## Notable Decisions

- **FILENAME convention**: Morpheus is instructed to start its response with `FILENAME: <word>` on line 1. The shell extracts this word with `grep -oP`, sanitises it, and falls back to `date +%Y%m%d` if parsing fails. This avoids a second API call.
- **`permissions: contents: write, models: read`**: Required for the workflow to push the commit back and to use the GitHub Models API with `GITHUB_TOKEN`.
- **Branch guard via `if:`**: Uses `github.event.repository.default_branch` to ensure the workflow only runs on the default branch (the `schedule` event always fires on the default branch anyway, but the guard is explicit as specified).
- **Multiline env passing**: Uses `<<HEREDOC` style (`<<ENDOFCONTEXT` / `ENDDREAM`) in `$GITHUB_ENV` to safely pass multi-line context and the LLM response between steps.
