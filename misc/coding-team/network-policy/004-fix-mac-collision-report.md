# Task 004 – Fix MAC address collision: Completion Report

## Summary

- Added `generateMAC()` helper in `internal/netutil/netutil.go` that generates 6 cryptographically-random bytes, sets the locally-administered bit (`0x02`), and clears the multicast bit (`0xfe` mask) on the first octet, then returns a `net.HardwareAddr`.
- In `ConfigureInterface`, inserted a call to `generateMAC()` and `netlink.LinkSetHardwareAddr(link, mac)` immediately after obtaining the link reference and **before** the `netlink.LinkSetUp(link)` call, per the constraint that MAC must be set before the interface is brought up.
- Added `"crypto/rand"` to the import block.

## Files changed

- `internal/netutil/netutil.go`

## Validation

`make test` and `make build` both passed with no errors or failures.

## Notable tradeoffs / risks

None. The change is minimal and contained to a single file. `crypto/rand.Read` is non-blocking on Linux (reads from `/dev/urandom`) so there is no risk of hanging in the container netns context.
