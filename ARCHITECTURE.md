# ARCHITECTURE.md

## Overview
The fancni project is a Go-based tool with a focus on network utilities. It leverages libraries for interacting with iptables and network namespaces. This appears to involve overlay capabilities (e.g., 'fan' configuration).

## Detected Stack
- **Language**: Go
- **Key Libraries**:
  - github.com/coreos/go-iptables
  - github.com/vishvananda/netlink and netns
- **Build System**:
  - `go.mod` for dependency management
  - `Makefile` for standardized commands (build, test, clean)

## Project Structure
- `cmd/fancni/`: Primary executable
- `internal/fan`: Fan (overlay) utilities
- `internal/iptables`: iptables helpers
- `internal/config`: Configuration handling
- `internal/netutil`: Network utility methods
- `misc`: Includes planning documents

## Linting, Testing & Build
- Run `make build`, `make test`

## Missing CI/CD
This project has no workflows.

---
