# Morpheus — consolidation-054251

_Generated: 2026-07-16T05:42:51Z_

# Idea: Dream File Consolidation

## Description

The `misc/nightly-dreams` directory currently contains a large number of Markdown files, many of which represent similar themes or duplicate ideas (e.g., numerous `consolidation-*.md`, `curation-*.md`, `prioritization-*.md`, and deduplication files). This proliferation makes it difficult to find the most valuable and actionable suggestions. The next step should be to **consolidate** these files: group together ideas that are thematically identical or substantially overlapping, and merge their content into a single, well-structured Markdown file for each core theme (e.g. `consolidation.md`, `curation.md`, `prioritization.md`, `deduplication.md`, etc.). Archive or remove the redundant, less clear, or lower-value files.

This will help maintainers and contributors focus on the most important, actionable, and non-redundant suggestions, and will make the process of nightly idea generation sustainable.

## Feasibility

- **Effort**: Moderate. Requires careful review and manual or scripted grouping/merging of dozens of files, but the process is mechanical and does not require deep code changes.
- **Impact**: High. Reduces clutter, improves maintainability, and increases the signal-to-noise ratio of the "dreams" directory, making it a much more valuable resource for planning future work. 
- **Dependencies**: None, can be started immediately. 

**Recommendation**: Begin consolidation with the most numerous and overlapping idea types (e.g. `consolidation-*`, `curation-*`, `prioritization-*`, `deduplication-*`). Add a brief rationale or summary section to each theme file describing the merged ideas.
