# Task: Create Nightly Dreams GitHub Actions Workflow

## Context
This is a Go project (fan overlay / CNI networking tool) with no existing CI/CD. We need a nightly GitHub Action that uses an LLM agent ("Morpheus") to analyze the project and suggest ideas.

## Objective
Create `.github/workflows/nightly-dreams.yml` — a nightly cron workflow that:
1. Checks out the repo.
2. Gathers context: project tree and git commits from the last 7 days.
3. Calls the GitHub Models API (chat completions) with that context to generate ranked ideas.
4. Writes the response to `misc/nightly-dreams/<one-word>.md` and commits it.

## Scope
- Single new file: `.github/workflows/nightly-dreams.yml`
- Also create the directory `misc/nightly-dreams/` with a `.gitkeep` so the path exists.

## How it should work

### Cron
- Schedule: nightly (e.g., `0 3 * * *`).
- Runs on default branch only.

### Gathering context
Use shell commands in the workflow to build a prompt payload:
- `find . -not -path './.git/*' -not -path './_output/*' | head -200` for project structure.
- `git log --since='7 days ago' --pretty=format:'%h %s' --stat` for recent commits (truncate to ~4000 chars to stay within token budget).
- Read key files if helpful (e.g., `go.mod`, `Makefile`) — keep it brief.

### Calling GitHub Models
- Use `curl` to call `https://models.github.ai/inference/chat/completions`.
- Auth: `Authorization: Bearer ${{ secrets.GITHUB_TOKEN }}`.
- Model: `openai/gpt-4.1` (available via GitHub Models).
- System prompt should identify the agent as "Morpheus" and instruct it to:
  - Analyze the project architecture and recent direction.
  - Suggest up to 10 concrete, ranked ideas for next steps.
  - Each idea: title, 2-3 sentence description, feasibility note.
  - Rank by impact and feasibility; discard extras beyond 10.
- The user message contains the gathered context (tree + commits + key files).
- Parse the response with `jq` to extract `.choices[0].message.content`.

### Writing output
- Generate a one-word description from the response (use jq or a second small LLM call, or just extract the first prominent keyword — keep it simple, e.g., use `date +%Y%m%d` as fallback if parsing is tricky). A simple approach: ask Morpheus in the same prompt to start the response with `FILENAME: <one-word>` on the first line, then strip it.
- Write to `misc/nightly-dreams/<one-word>.md`.
- Commit and push with a bot-style commit message like `morpheus: nightly dreams — <word>`.
- Use `git config user.name "Morpheus"` and a noreply email.
- Add `[skip ci]` to the commit message to avoid recursive triggers.

## Non-goals
- No PR creation.
- No notifications.
- No manual trigger / workflow_dispatch.
- No retry logic or error alerting.

## Constraints
- No extra secrets — only `GITHUB_TOKEN`.
- Keep the workflow self-contained (no external actions beyond `actions/checkout`).
- The workflow must handle the case where there are no commits in the last 7 days gracefully (just skip and exit 0).
