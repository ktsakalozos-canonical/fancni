# Task 004 Completion Report: CNI Binary — ADD/DEL/CHECK/VERSION

## Summary

- Added `internal/fan/fanctl.go` — `EnsureBridge` wrapping the only `exec.Command` in the codebase (`fanctl up -o … -u …`), with a clear error when fanctl is not in PATH.
- Added `internal/cni/plugin.go` — `Plugin` struct with `HandleAdd`, `HandleDel`, `HandleCheck`, `HandleVersion`; CNI 1.0.0 result/error JSON structures; random veth naming via `crypto/rand`.
- Replaced `cmd/fancni/main.go` stub — reads config from stdin, detects host IP via UDP dial trick, wires together IPAM + Plugin, dispatches `CNI_COMMAND`, writes CNI error JSON to stdout on failure and exits 1.
- Added `internal/cni/plugin_test.go` — tests for `HandleVersion` output, `cniResult` JSON shape (no `version` field in ips per CNI 1.0.0), `randomVethName` uniqueness, and `HandleDel` idempotency.

## Files Changed

- `internal/fan/fanctl.go` (new)
- `internal/cni/plugin.go` (new)
- `internal/cni/plugin_test.go` (new)
- `cmd/fancni/main.go` (replaced stub)

## Notable Tradeoffs / Risks

- **HandleDel skips veth cleanup intentionally** — per the task brief, both veth ends are cleaned up automatically when the container netns is destroyed by the runtime. This matches Flannel/Calico behaviour.
- **VERSION command skips stdin read** — `CNI_COMMAND=VERSION` is dispatched before `config.Parse(os.Stdin)` because the runtime may not send a config body for VERSION. This is consistent with the CNI spec.
- **Log file fallback** — if `/var/log/fancni.log` is not writable (e.g. in test environments), logging falls back to stderr silently so the binary still functions.
