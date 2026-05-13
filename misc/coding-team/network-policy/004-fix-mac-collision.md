# Task: Fix MAC address collision in veth creation

## Context

On kernel 6.8 (Ubuntu 24.04), creating multiple veth pairs results in duplicate MAC addresses on the container-side interfaces. This prevents same-node pod-to-pod communication because the bridge can't distinguish between pods at L2 (ARP responses get misdelivered).

Evidence: 3 out of 4 pods on the same node got MAC `52:78:2c:a2:b3:45`.

## Objective

Assign a unique, random, locally-administered MAC address to each container-side veth interface.

## Scope

Modify `internal/netutil/netutil.go`:

Add a helper function to generate a random MAC with the locally-administered bit set:
```go
func generateMAC() (net.HardwareAddr, error) {
    buf := make([]byte, 6)
    if _, err := rand.Read(buf); err != nil {
        return nil, fmt.Errorf("generating random MAC: %w", err)
    }
    // Set locally-administered bit, clear multicast bit
    buf[0] = (buf[0] | 0x02) & 0xfe
    return net.HardwareAddr(buf), nil
}
```

Then in `ConfigureInterface`, after getting the link and before bringing it up, set the MAC:
```go
mac, err := generateMAC()
if err != nil {
    return err
}
if err := netlink.LinkSetHardwareAddr(link, mac); err != nil {
    return fmt.Errorf("setting MAC on %q: %w", ifName, err)
}
```

Add `"crypto/rand"` to the imports.

## Why ConfigureInterface (not CreateVethPair)

`ConfigureInterface` already runs inside the container netns. Setting the MAC there is the safest because:
1. No collision with other interfaces in the host namespace
2. The link is already moved to the target netns
3. We have a clear reference to the correct interface

## Constraints

- Use `crypto/rand` (not `math/rand`) for MAC generation
- Set locally-administered bit (0x02 on first octet) and clear multicast bit (0xfe mask on first octet)
- The MAC must be set BEFORE bringing the link up (order matters: set MAC, then LinkSetUp)

## Non-goals

- Do not change CreateVethPair or the veth name generation
- Do not add a unit test for generateMAC (it's trivial)
