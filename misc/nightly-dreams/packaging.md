# Morpheus — packaging

_Generated: 2026-05-15T11:31:42Z_

# Packaging Optimization

## Title
Enhanced Packaging and Distribution for Rocks

## Description
With the recent migration from Docker to Rocks (as evidenced by the removal of Docker artifacts and the addition of Rockcraft configuration), the project now builds and distributes its software as Linux containers using Rocks. The next logical step is to optimize the packaging and distribution process:
- Automate version tagging and release notes generation for new rock builds.
- Integrate checks for reproducible builds and image size optimization in the CI workflows.
- Explore publishing rocks to multiple registries (e.g., GHCR, Canonical’s own registry) for broader accessibility.
- Provide clear user documentation for installing and running the rocks, replacing outdated Docker instructions.

This would help streamline the release process, improve reliability, and ensure easy access for end-users and downstream teams.

## Feasibility Note
This is highly feasible. The project already uses GitHub Actions for CI/CD, and Rockcraft for packaging. Extending the workflows and documentation requires low to moderate effort, mostly in YAML and Markdown, with some scripting. The benefit is significant for both maintainers and users.
