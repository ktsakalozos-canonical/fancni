// Package ipam provides IP address management for fancni containers.
package ipam

import "net"

// IPAM manages the allocation and release of IP addresses for containers.
type IPAM interface {
	// Allocate assigns an IP to the given containerID, returning the IP.
	// If the containerID already has an allocation, the existing IP is returned (idempotent).
	Allocate(containerID string) (net.IP, error)

	// Lookup returns the IP currently assigned to containerID.
	// Returns (nil, false, nil) if the containerID has no allocation.
	Lookup(containerID string) (net.IP, bool, error)

	// Free releases the IP assigned to containerID.
	// Returns the freed IP and true if found, or (nil, false, nil) if not found.
	Free(containerID string) (net.IP, bool, error)
}
