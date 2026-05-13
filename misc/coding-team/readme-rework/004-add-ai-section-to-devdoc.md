# Task: Add AI-Assisted Development section to docs/development.md

## Context
The project uses OpenCode with configured AI agents and two nightly GitHub Actions workflows. Contributors should know these exist.

## Objective
Add an "AI-Assisted Development" section to `docs/development.md` (after the existing Makefile section). Awareness-level only — what exists and where to find it.

## Content to include

### Agents (`.opencode/agents/`)
The project uses OpenCode with 5 subagents:
- **architect** — orchestrates work, writes task briefs, delegates to other agents
- **developer** — implements tasks
- **code-reviewer** / **code-reviewerer** — dual independent code reviews
- **repo-scout** — researches the repo, maintains `ARCHITECTURE.md`

Agent prompts live in `.opencode/agents/`. Each agent reads a corresponding knowledge file from `misc/knowledge/` at session start for project-specific context.

### Nightly Knowledge Refinement (`.github/workflows/nightly-knowledge.yml`)
A scheduled workflow that uses the GitHub Models API to generate/refine per-agent knowledge files under `misc/knowledge/`. Runs daily; only triggers if there have been commits in the last 7 days.

### Morpheus (`.github/workflows/nightly-dreams.yml`)
A nightly "Morpheus" agent that analyzes the project and suggests improvement ideas. Ideas are stored as individual files under `misc/nightly-dreams/`.

## Constraints
- Keep it concise — awareness level, not operational instructions
- Don't restructure existing content in the file
- Don't add instructions on how to create new agents or modify workflows

## Non-goals
- No separate doc file — this goes in the existing `docs/development.md`
- No detailed workflow configuration docs
