# Development

## Prerequisites

- Go 1.24+

## Makefile

The Makefile is the canonical entry point for all development commands.

```bash
# Build the CNI binary
make build

# Run tests
make test

# Build rock image
make rock-build

# Validate Helm chart
make helm-lint
make helm-template
```

Internally, the project uses standard `go build ./...` and `go test ./...`; no special toolchain beyond Go 1.24+ is required.

## AI-Assisted Development

The project uses [OpenCode](https://opencode.ai) with a set of configured subagents and nightly GitHub Actions workflows.

### Agents (`.opencode/agents/`)

Five subagents are configured:

- **architect** — orchestrates work, writes task briefs, and delegates to other agents
- **developer** — implements tasks defined by the architect
- **code-reviewer** / **code-reviewerer** — two independent code reviewers that evaluate implementations
- **repo-scout** — researches the repository and maintains `ARCHITECTURE.md`

Agent prompts live in `.opencode/agents/`. Each agent reads a corresponding knowledge file from `misc/knowledge/` at session start for project-specific context.

### Nightly Knowledge Refinement (`.github/workflows/nightly-knowledge.yml`)

A scheduled workflow that uses the GitHub Models API to generate and refine per-agent knowledge files under `misc/knowledge/`. Runs daily; only triggers if there have been commits in the last 7 days.

### Morpheus (`.github/workflows/nightly-dreams.yml`)

A nightly agent that analyzes the project and suggests improvement ideas. Ideas are stored as individual files under `misc/nightly-dreams/`.
