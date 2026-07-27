# Morpheus — pruning

_Generated: 2026-07-27T06:32:42Z_

# Idea: Pruning Old Nightly Dream Files

## Description

The `misc/nightly-dreams` directory is accumulating a large number of nightly idea files, including several with similar themes (e.g., consolidation, curation, observability) and potential duplicates over time. This can create clutter and make it harder for contributors to find the most relevant or actionable ideas.

Introduce a scheduled or manual pruning process:
- **Keep only the 10 most recent and/or most relevant/unique nightly dream files.**
- Consider a simple script (e.g., in Python, Bash, or Go) that runs as a GitHub Action or locally, which:
  - Sorts files by creation/modification time or by a relevance/feasibility label.
  - Deletes older files or those marked as superseded/obsolete.
  - Optionally, consolidates similar ideas into a single summary file before deletion.

This will help keep the idea pool actionable, relevant, and manageable for maintainers and contributors.

## Feasibility

This is highly feasible. The project already uses GitHub Actions and shell scripts; adding a simple cleanup script is straightforward. No code changes are needed in the core Go project, only in housekeeping. The main decision will be the policy for relevance/selection, which can evolve over time.
