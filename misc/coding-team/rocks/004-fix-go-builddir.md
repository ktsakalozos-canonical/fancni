# Task: Remove invalid go-builddir field

## Context
`rockcraft pack` fails because `go-builddir` is not a valid field for the Go plugin. The Go plugin only accepts `go-buildtags` and `go-generate`. It always runs `go install ./...` which builds all packages.

## Objective
Remove the `go-builddir: cmd/fancni` line from the `fancni-binary` part in `rockcraft.yaml`.

## Notes
- The Go plugin runs `go install ./...` by default, which will build `cmd/fancni` and place the binary in `bin/fancni` — so the existing `organize` mapping is still correct.
- No other changes needed.
