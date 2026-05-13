# Task: Update agent prompts to read knowledge files

## Context
A nightly workflow now generates `misc/knowledge/<agent>-knowledge.md` for each of the 5 agents. These files contain role-specific project knowledge to give agents a head-start.

## Objective
Update each agent prompt in `.opencode/agents/` to instruct the agent to read its knowledge file at the start of a session.

## Scope
Files to modify:
- `.opencode/agents/architect.md`
- `.opencode/agents/developer.md`
- `.opencode/agents/code-reviewer.md`
- `.opencode/agents/code-reviewerer.md`
- `.opencode/agents/repo-scout.md`

For each file, add a short section (2-4 lines) near the top of the prompt body (after the YAML frontmatter, before or within the "Project/stack awareness" or equivalent section) instructing the agent to:
1. Read `misc/knowledge/<agent>-knowledge.md` at the start of each session (if it exists)
2. Treat it as supplementary context — useful but not authoritative (the agent should still verify if something seems off)

## Non-goals
- Don't restructure or rewrite the existing prompts
- Don't make the knowledge file mandatory (use "if it exists" language)
- Don't change agent behavior, permissions, or mode

## Constraints
- Keep additions minimal — 2-4 lines per file
- Place the instruction where it flows naturally with existing prompt structure
