# Task: Fix YAML indentation in nightly-knowledge workflow

## Context
The multiline bash string assignments for `USER_MESSAGE` in `.github/workflows/nightly-knowledge.yml` have content lines at column 0, which breaks YAML literal block scalar parsing. The `run: |` block requires all non-empty lines to be indented at least to the block's base indentation (10 spaces).

## Objective
Fix the YAML indentation for all 5 "Generate X knowledge" steps so the file parses correctly.

## Scope
File: `.github/workflows/nightly-knowledge.yml`

In each of the 5 generate steps, the pattern:
```
            USER_MESSAGE="${PROJECT_CONTEXT}

=== Prior Knowledge ===
${PRIOR_KNOWLEDGE}"
```

Must become (indented to at least 10 spaces — matching the block's base indentation):
```
            USER_MESSAGE="${PROJECT_CONTEXT}

          === Prior Knowledge ===
          ${PRIOR_KNOWLEDGE}"
```

This affects all 5 instances (architect, developer, code-reviewer, code-reviewerer, repo-scout).

Note: YAML will strip the 10-space base indentation, so the actual bash string value remains `\n=== Prior Knowledge ===\n${PRIOR_KNOWLEDGE}` — no extra whitespace in the content.

## Non-goals
- Don't change anything else in the file

## Acceptance criteria
- The workflow YAML parses without errors (validate with a YAML linter or `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/nightly-knowledge.yml'))"`)
