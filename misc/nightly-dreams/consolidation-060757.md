# Morpheus — consolidation-060757

_Generated: 2026-07-26T06:07:57Z_

# Idea: Nightly Dreams File Consolidation

## Description

The `misc/nightly-dreams` directory is accumulating a large and growing number of idea files, including multiple with similar themes (e.g., "consolidation-054450.md", "consolidation-055726.md", etc.), and periodic curation files. To keep the project maintainable and make it easier for contributors to review and act on proposed enhancements, introduce an automated or semi-automated consolidation process. This process should:

- Regularly scan the directory for overlapping, duplicate, or outdated ideas.
- Merge similar ideas into a single, clearly titled file with cross-references to original contributors if appropriate.
- Archive or remove ideas that are obsolete, already implemented, or superseded by newer suggestions.
- Generate a brief summary/index file listing active, pending, and recently consolidated ideas for quick navigation.

This can be implemented as a standalone script (e.g., in Go or Python) and/or as a GitHub Action, possibly integrated into the existing nightly workflow.

## Feasibility

**High.**  
The project already has experience with workflow automation and curation (as evidenced by existing workflow YAMLs and curation files). The main work involves scripting file scanning, basic text processing, and file operations—all straightforward tasks. This will greatly enhance the long-term maintainability and efficiency of the "nightly dreams" ideation process.
