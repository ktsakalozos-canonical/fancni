# Morpheus — deduplication

_Generated: 2026-05-21T14:19:37Z_

# Idea: Nightly Dreams Deduplication

## Description

The `misc/nightly-dreams` directory has accumulated over 10 idea files, some of which overlap in scope (e.g., `maintainability.md`, `cleanup.md`, `modularity.md`) or are now less relevant given recent development directions (such as the shift to rocks packaging). To ensure the nightly-dreams serve as actionable and inspiring next steps, introduce a nightly deduplication and pruning script or process. This process should:

- Automatically scan `misc/nightly-dreams` for more than 10 idea files.
- Group ideas by similarity or theme, surfacing only the most actionable, unique, and relevant 10.
- Delete or archive less feasible or redundant suggestions.
- Optionally, generate a short summary report to justify why each idea was kept or pruned.

This will keep the nightly-dreams focused, fresh, and maximally useful for the team.

## Feasibility

**High.**  
The process can be implemented as a shell or Python script integrated into the existing GitHub Actions (see `.github/workflows/nightly-dreams.yml`). It builds on existing nightly automation patterns, requires no changes to core project code, and will streamline ideation for all contributors.
