# Morpheus — curation

_Generated: 2026-05-22T12:08:03Z_

# Idea: Nightly Dream Curation System

## Description

The `misc/nightly-dreams` directory has accumulated a large number of improvement suggestions, with many topics now covered (e.g., packaging, performance, testing, analytics, maintainability). To prevent idea overload and ensure the highest-impact improvements are prioritized, implement a **curation system** for nightly dreams. This process would periodically review, score, and select the top 10 most relevant, actionable, and feasible ideas, archiving or discarding less relevant or outdated ones. This ensures maintainers and contributors are always focused on the most valuable next steps, and prevents the directory from growing unwieldy.

Curation could be manual (scheduled review sessions), automated (a bot that scores and prunes ideas based on criteria), or a hybrid (suggestions are flagged for review after a threshold is reached).

## Feasibility Note

This is a light-to-moderate effort, depending on implementation. A manual curation policy can be adopted immediately. Automating curation with a simple script or GitHub Action is straightforward, as all ideas are in one directory and follow a clear file format. This idea improves project focus and maintainability without requiring major architectural changes.
