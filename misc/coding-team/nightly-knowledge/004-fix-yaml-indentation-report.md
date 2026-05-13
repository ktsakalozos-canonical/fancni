# Report: Fix YAML indentation in nightly-knowledge workflow

## Summary
- Fixed 5 instances of unindented lines inside `run: |` blocks in `.github/workflows/nightly-knowledge.yml`
- The `=== Prior Knowledge ===` and `${PRIOR_KNOWLEDGE}"` lines were at column 0, breaking YAML literal block scalar parsing; they are now indented to 10 spaces (matching the block's base indentation)
- YAML validation passes: `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/nightly-knowledge.yml'))"`

## Files changed
- `.github/workflows/nightly-knowledge.yml`

## Notable tradeoffs or risks
- None. The fix is purely whitespace. YAML strips the base indentation, so the actual bash string values are unchanged.
