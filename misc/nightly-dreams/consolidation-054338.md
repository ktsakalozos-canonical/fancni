# Morpheus — consolidation-054338

_Generated: 2026-07-30T05:43:38Z_

# Idea: Nightly Dream File Consolidation

## Description

The `misc/nightly-dreams` directory is accumulating a growing number of files, some of which are single-idea records and others are time-stamped or batch "consolidation", "selection", or "pruning" files. Over time, this can make it difficult to track which ideas are current, which have been superseded or selected, and which should be kept or removed. 

**Proposal:**  
Introduce a nightly or weekly automated consolidation process that reviews all files in `misc/nightly-dreams`, selects the 10 most relevant and feasible ideas, and moves the rest to an archival subdirectory (e.g., `misc/nightly-dreams/archive/`). This process should also update a single summary file (e.g., `misc/nightly-dreams/README.md`) that lists the current "active" top 10 ideas with links to their markdown files and short rationales for their selection.

This will keep the directory organized, make it easier for maintainers to review ongoing proposals, and ensure that only the best ideas are visible for further consideration and implementation.

## Feasibility

- **Technical:** High. Implementing a cron-triggered GitHub Action (or reusing the existing nightly jobs) to scan, select, archive, and update a README is straightforward and can be written in bash, Python, or Go.
- **Process:** Moderate. Requires a clear, automated selection heuristic (e.g., most recent, by filename, or by parsing feasibility/interest comments).
- **Overall:** This is a maintainable improvement that will prevent idea sprawl and confusion, with minimal ongoing overhead.
