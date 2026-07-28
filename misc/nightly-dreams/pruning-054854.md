# Morpheus — pruning-054854

_Generated: 2026-07-28T05:48:54Z_

# Pruning Old Nightly Dreams

## Description

The `misc/nightly-dreams` directory has accumulated a large number of files representing past ideas and nightly suggestions. As the project evolves, these files can become outdated, irrelevant, or even conflicting with more recent directions. To keep the idea pool actionable and relevant, implement a regular pruning process: limit the number of active nightly-dreams files to 10. Each time a new idea is generated, review all existing ideas and discard the least interesting or feasible ones, ensuring only the best 10 remain. This will help maintain focus and prevent idea fatigue.

## Feasibility

This is immediately feasible and mostly a process/policy change. The implementation can be automated via the existing nightly GitHub workflows (see `.github/workflows/nightly-dreams.yml`), or simply reinforced through clear contributor documentation. This approach keeps the idea pool fresh and impactful.
