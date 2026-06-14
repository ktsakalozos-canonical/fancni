# Morpheus — consolidation-072504

_Generated: 2026-06-14T07:25:04Z_

# Nightly Idea Consolidation

## Description

The `misc/nightly-dreams` directory currently contains a large number of suggestion files, many of which are one-off or duplicate "consolidation-*" and "curation-*" files, plus a variety of themed improvement ideas (e.g., `usability.md`, `testing.md`, `security.md`). This level of fragmentation makes it difficult to survey the project's improvement history, prioritize actionable ideas, and prevent duplication of effort.

**Proposal:**  
Introduce a weekly or bi-weekly consolidation process that reviews all files in `misc/nightly-dreams`, merges similar/duplicate content, and selects the 10 most relevant, actionable, and feasible ideas. This process should delete outdated or superseded files and update a `consolidation.md` file that lists the current top ideas, their status, and which ones are being actively pursued.

## Feasibility

**High.**  
This process can be implemented as a lightweight scripted workflow (possibly using a Makefile target or a simple Go script), or as a documented manual process. It leverages existing suggestion files but makes the directory more maintainable and valuable for the team. This consolidation also improves onboarding and quarterly planning, and can be done without code changes to the main product.
