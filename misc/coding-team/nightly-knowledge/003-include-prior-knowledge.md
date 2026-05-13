# Task: Include prior knowledge in nightly generation prompts

## Context
The nightly-knowledge workflow (`.github/workflows/nightly-knowledge.yml`) generates knowledge files for each agent. Currently, each LLM call only receives project context. The user wants the LLM to also receive the existing knowledge file as input so it can refine/build on prior output.

## Objective
For each of the 5 agent LLM calls, read the existing `misc/knowledge/<agent>-knowledge.md` file (if it exists) and include it in the prompt as prior knowledge.

## Scope
File to modify: `.github/workflows/nightly-knowledge.yml`

For each of the 5 "Generate X knowledge" steps:
1. Before calling the API, read the existing knowledge file into a variable (handle the case where it doesn't exist yet)
2. Append the prior knowledge to the user message (or add it as a separate section in the user content), clearly labeled (e.g., "=== Prior Knowledge ===")
3. Add a line to the system prompt telling the LLM it can use/refine/update the prior knowledge, keeping what's still accurate and updating what's changed

## Non-goals
- Don't restructure the workflow
- Don't change the system prompt focus areas or output rules
- Don't change the write/commit steps

## Constraints
- Keep the addition minimal per step — read file + append to context
- Handle missing file gracefully (first run won't have prior knowledge)
