# Morpheus — prioritization-064305

_Generated: 2026-07-05T06:43:05Z_

# Prioritization of Nightly Dreams

## Description

The `misc/nightly-dreams` directory currently contains a large number of idea files, many of which focus on overlapping themes such as consolidation, deduplication, and curation. With over 40 files present, it is challenging for maintainers to identify the most relevant, actionable, and non-redundant suggestions. The next step should be to implement a systematic prioritization and cleanup process:

- **Review all ideas in `misc/nightly-dreams`.**
- **Select the 10 most interesting and feasible ideas**, ensuring diversity (e.g., covering testing, observability, security, performance, usability, etc.).
- **Delete or archive the remaining files** to avoid clutter and focus the team's efforts on the most impactful work.
- **Establish a policy or lightweight workflow** for future idea submissions, e.g., ideas must reference existing files and justify their addition if the directory is full.

This will help keep the project's improvement pipeline actionable, reduce cognitive overload, and make it easier for contributors to align on priorities.

## Feasibility Note

This can be accomplished with a one-time focused review session (possibly with input from core maintainers or an automated script), followed by periodic reviews if the number of idea files grows again. As this is a process and documentation improvement, it requires no code changes and is low-risk and highly feasible.
