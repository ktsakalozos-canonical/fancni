# Morpheus — deduplication-070250

_Generated: 2026-06-20T07:02:50Z_

# Nightly Dream: Deduplication of Nightly Dreams

## Description

The `misc/nightly-dreams` directory has accumulated a large number of idea files, many with similar or overlapping content. There are multiple "consolidation", "curation", and "prioritization" variants, as well as older residual files. This proliferation makes it difficult to identify the most relevant or actionable ideas, and risks valuable insights being lost in the noise.

**Proposal:**  
Introduce a regular deduplication process for the `nightly-dreams` directory. The process would:
- Periodically review all files and merge similar ideas into a single, canonical file.
- Remove outdated or redundant entries.
- Maintain a cap (e.g., 10 files) to ensure only current, high-priority ideas are kept.
- Optionally, add a changelog or archive for merged/discarded ideas for traceability.

This will keep the idea pool actionable and focused, and make it easier for contributors to discover and implement the best suggestions.

## Feasibility

Highly feasible. This is primarily an organizational and maintenance task, which could be automated with a simple script or handled as part of the Morpheus agent's workflow. The main challenge is semantic similarity detection, which can start with simple heuristics (e.g., filename and keyword matching) and become more sophisticated over time. The process can also include a manual review step if needed.
