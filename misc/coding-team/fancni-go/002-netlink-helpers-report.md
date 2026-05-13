# Task 002: Netlink Helpers — Completion Report

## Summary

- Added `github.com/vishvananda/netlink`, `github.com/vishvananda/netns`, and `github.com/coreos/go-iptables` as dependencies via `go get`.
- Created `internal/netutil/netutil.go` with all required netlink-based functions (`CreateVethPair`, `AttachToBridge`, `MoveToNetNS`, `ConfigureInterface`, `DeleteLink`, `GetLinkMAC`, `LinkExists`, `VerifyInterfaceConfig`). No `exec.Command` calls — all via netlink. Namespace-switching functions lock the OS thread and always restore the original namespace via defer.
- Created `internal/iptables/iptables.go` with `EnsureForwardRules` and `EnsureMasqueradeRule`, both using `AppendUnique` for idempotency.
- Fixed `internal/config/config.go` to reject `underlayPrefix > 32` (changed bounds check to `< 1 || > 32`).

## Files Changed

- `go.mod` / `go.sum` — new dependencies
- `internal/config/config.go` — tightened `underlayPrefix` bounds check
- `internal/config/config_test.go` — added `TestParse_UnderlayPrefixTooLarge`
- `internal/netutil/netutil.go` — new file
- `internal/netutil/netutil_test.go` — new file (root-required tests skipped; `LinkExists` tested without root)
- `internal/iptables/iptables.go` — new file
- `internal/iptables/iptables_test.go` — new file (all tests skip, require root)

## Notable Tradeoffs / Risks

- `isExist` uses a string comparison on `err.Error()` as a fallback alongside `LinkNotFoundError` type assertion; the netlink library does not always wrap `syscall.EEXIST` in a typed error. This is acceptable for now and will be exercised properly in integration tests with root.
- Namespace-switching functions (`ConfigureInterface`, `GetLinkMAC`, `VerifyInterfaceConfig`) call `runtime.LockOSThread` to prevent the goroutine scheduler from moving the goroutine to a thread that hasn't had its namespace changed — standard practice for netns operations.
