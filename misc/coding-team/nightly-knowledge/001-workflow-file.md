# Task: Create nightly-knowledge GitHub Actions workflow

## Context
The project already has `.github/workflows/nightly-dreams.yml` which uses the GitHub Models API (`https://models.github.ai/inference/chat/completions`) with `$GITHUB_TOKEN`. Follow the same pattern.

## Objective
Create `.github/workflows/nightly-knowledge.yml` — a single workflow that generates 5 agent-specific knowledge files under `misc/knowledge/`.

## Scope

### Workflow structure
- Trigger: `workflow_dispatch` + `schedule: cron '0 4 * * *'` (one hour after Morpheus)
- Single job, condition: only on default branch
- Permissions: `contents: write`, `models: read`
- Steps:
  1. Checkout (fetch-depth: 0)
  2. Check for recent commits (last 7 days) — skip if none
  3. Gather context (same as Morpheus: tree, recent commits, go.mod, Makefile, plus sample a few key source files from `internal/` and `cmd/`)
  4. One step per agent (5 steps) calling GitHub Models API with a role-specific system prompt
  5. Write files step — write each response to `misc/knowledge/<agent>-knowledge.md`, truncate to 500 lines
  6. Commit and push (only if files changed)

### LLM call details
- Model: `openai/gpt-4.1`
- Temperature: 0.3 (lower than Morpheus — we want factual, not creative)
- max_tokens: 4096
- Each call gets the same project context as user message, but a different system prompt

### System prompts (role-specific focus)
Each prompt must instruct the LLM to:
- Output pure markdown, no preamble
- Stay under 500 lines
- Be specific and actionable, not generic

Role-specific angles:
- **architect**: System boundaries and module responsibilities, architectural decisions evident in code, dependency graph, incomplete/in-progress work (from recent commits and TODOs), areas of technical debt
- **developer**: Code patterns and idioms used (error handling, logging, config), naming conventions, testing patterns and framework usage, build/run commands, common gotchas visible in code
- **code-reviewer**: Areas with weak/missing test coverage, security-sensitive code paths, code smells or inconsistencies currently present, error handling gaps
- **code-reviewerer**: Same focus as code-reviewer but emphasize edge cases, race conditions, resource leaks, and maintainability concerns
- **repo-scout**: Dependency list with versions, build/deploy tooling state, CI workflow inventory, project structure map with hotspots

### Commit step
- git config user: "Knowledge Bot" / "knowledge-bot@users.noreply.github.com"
- `git add misc/knowledge/`
- Only commit if `git diff --cached --quiet` fails (i.e., there are changes)
- Commit message: `knowledge-bot: update agent knowledge files [skip ci]`

## Non-goals
- No web searches or external API calls beyond GitHub Models
- No accumulation across runs — each run overwrites
- Don't touch agent prompt files (that's task 002)

## Constraints
- Knowledge files MUST NOT exceed 500 lines. Enforce with `head -500` after writing.
- Use the same API pattern as nightly-dreams.yml (curl to models.github.ai, jq to parse)
- Keep the workflow DRY where reasonable (shared context gathering) but don't over-abstract with reusable actions
