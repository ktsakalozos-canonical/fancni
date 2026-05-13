# Task: Create docs/development.md

## Context
The README currently has a "Development" section with build/test/lint commands. We're splitting it out for a cleaner README.

## Objective
Create `docs/development.md` containing development workflow information for contributors.

## Content to include
- Go version requirement (1.24+)
- Makefile targets: `make build`, `make test`, `make docker-build`, `make helm-lint`, `make helm-template`
- Brief note that the project uses standard `go build ./...` / `go test ./...`
- Mention the Makefile as the canonical entry point for all dev commands

## Non-goals
- No CI/CD setup instructions
- No deployment/release process docs
- Don't duplicate architecture info
