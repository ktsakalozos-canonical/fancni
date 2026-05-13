# Task 003: File-Based IPAM

## Context

Phases 1-2 built fan address math and netlink helpers. This task implements the IP address management (IPAM) component. The CNI binary (Phase 4) will call IPAM to allocate/free IPs per container.

## Objective

Create `internal/ipam/` with a file-backed IPAM that maps container IDs to IPs, guarded by `syscall.Flock` for concurrency safety.

## Scope

### internal/ipam/ipam.go — Interface

```go
type IPAM interface {
    Allocate(containerID string) (net.IP, error)
    Lookup(containerID string) (net.IP, bool, error)
    Free(containerID string) (net.IP, bool, error)
}
```

Note: `Free` takes `containerID` (not IP). This is simpler than the PoC which took IP and had to reverse-lookup.

### internal/ipam/file_ipam.go — Implementation

**Constructor:**
```go
func NewFileIPAM(dataDir, podCIDR string) *FileIPAM
```
- `dataDir`: directory for state files (e.g., `/var/lib/cni/fancni/`)
- `podCIDR`: the node's pod subnet (e.g., `240.3.4.0/24`)
- Allocation file: `<dataDir>/ipam.json`
- Lock file: `<dataDir>/ipam.lock`

**State file format** (`ipam.json`):
```json
{
  "containerID1": "240.3.4.2",
  "containerID2": "240.3.4.3"
}
```

**Allocate(containerID string) (net.IP, error)**
1. Acquire flock (LOCK_EX) on lock file, defer release
2. Read allocation file (create if doesn't exist → empty map)
3. If containerID already has an allocation, return it (idempotent for retries)
4. Build set of allocated IPs
5. Iterate from `.2` to `.254` in the CIDR, find first free
6. Write updated map back to file
7. Return allocated IP

**Lookup(containerID string) (net.IP, bool, error)**
1. Acquire flock, defer release
2. Read allocation file
3. Return IP if found, (nil, false, nil) if not

**Free(containerID string) (net.IP, bool, error)**
1. Acquire flock, defer release
2. Read allocation file
3. If containerID not found, return (nil, false, nil)
4. Delete entry, write back
5. Return freed IP

**Locking details:**
- Open lock file with `os.OpenFile(path, O_CREATE|O_RDWR, 0644)`
- `syscall.Flock(fd, LOCK_EX)` — blocking (no retry loop needed, the kernel blocks)
- Defer `syscall.Flock(fd, LOCK_UN)` and `file.Close()`
- Create `dataDir` with `os.MkdirAll` if it doesn't exist (in constructor or lazily)

### Unit tests — internal/ipam/file_ipam_test.go

Use `t.TempDir()` for isolated test directories. These tests do NOT require root.

Test cases:
- Allocate first IP → should get `.2`
- Allocate second IP → should get `.3`
- Allocate same containerID twice → idempotent, returns same IP
- Lookup existing → found
- Lookup non-existing → not found
- Free existing → returns freed IP, subsequent lookup returns not found
- Free non-existing → returns (nil, false, nil)
- Allocate after free → reuses freed IP
- Exhaust all IPs (`.2` through `.254`) → 253 allocations succeed, 254th fails

## Non-goals

- No netlink, no iptables, no fanctl
- No CNI command handling
- No DHCP or remote IPAM backends

## Constraints

- Use `syscall.Flock` (blocking mode), not a retry loop
- Allocation file must be valid JSON at all times (atomic write: write to temp file, rename)
- `dataDir` created automatically if missing
