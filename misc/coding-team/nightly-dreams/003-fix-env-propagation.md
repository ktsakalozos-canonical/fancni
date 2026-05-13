# Task: Fix DREAM_WORD env propagation on collision

## Context
In `.github/workflows/nightly-dreams.yml`, the "Write dream file" step (line 134-136) updates `DREAM_WORD` with a timestamp suffix when a filename collision occurs. However, this is a local shell variable change — it does NOT propagate to the next step ("Commit and push") because GitHub Actions steps run in separate shell processes. The `GITHUB_ENV` value from the earlier step is still the old value.

## Objective
When a collision is detected and `DREAM_WORD` is updated, write the new value to `$GITHUB_ENV` so the "Commit and push" step picks it up.

## Change
In the "Write dream file" step, after updating `DREAM_WORD` on collision (line 135), add:
```
echo "DREAM_WORD=$DREAM_WORD" >> "$GITHUB_ENV"
```

Also update the `git add` line in "Commit and push" to use the output file path more robustly. Actually the simplest approach: merge the write + commit into a single step, or pass the final filename via GITHUB_ENV. The GITHUB_ENV approach is simpler — just add the echo line.

## Non-goals
- No other changes.
