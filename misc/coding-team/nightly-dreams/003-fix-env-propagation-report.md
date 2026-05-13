# Report: Fix DREAM_WORD env propagation on collision

## Summary
- Added `echo "DREAM_WORD=$DREAM_WORD" >> "$GITHUB_ENV"` after the collision-handling update in the "Write dream file" step.
- This ensures the updated `DREAM_WORD` (with timestamp suffix) is propagated to the subsequent "Commit and push" step via `GITHUB_ENV`, so `git add` and `git commit` reference the correct filename.

## Files changed
- `.github/workflows/nightly-dreams.yml`

## Notable tradeoffs or risks
- None. The change is minimal and targeted exactly as specified in the Task Brief.
