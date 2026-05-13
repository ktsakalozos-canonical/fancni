// Package fan provides pure functions for computing Fan Networking addresses.
//
// Fan Networking maps an underlay network (node IPs, typically /16) to an
// overlay network (pod IPs, /8) deterministically. Given underlay node IP
// 172.16.3.4 and overlay 240.0.0.0/8, the node's pod subnet is 240.3.4.0/24
// (overlay first octet + underlay 3rd and 4th octets).
package fan

import (
	"fmt"
	"net"
)

// ComputeSubnet returns the /24 pod subnet for a node given the overlay network
// CIDR and the node's host IP.
//
// Example: overlayNetwork="240.0.0.0/8", hostIP=172.16.3.4 → "240.3.4.0/24"
func ComputeSubnet(overlayNetwork string, hostIP net.IP) (string, error) {
	overlayFirst, err := validateOverlay(overlayNetwork)
	if err != nil {
		return "", err
	}

	host4 := to4(hostIP)
	if host4 == nil {
		return "", fmt.Errorf("hostIP must be an IPv4 address, got: %v", hostIP)
	}

	return fmt.Sprintf("%d.%d.%d.0/24", overlayFirst, host4[2], host4[3]), nil
}

// ComputeGateway returns the gateway IP (.1) of the computed fan subnet.
//
// Example: overlayNetwork="240.0.0.0/8", hostIP=172.16.3.4 → 240.3.4.1
func ComputeGateway(overlayNetwork string, hostIP net.IP) (net.IP, error) {
	overlayFirst, err := validateOverlay(overlayNetwork)
	if err != nil {
		return nil, err
	}

	host4 := to4(hostIP)
	if host4 == nil {
		return nil, fmt.Errorf("hostIP must be an IPv4 address, got: %v", hostIP)
	}

	gw := net.ParseIP(fmt.Sprintf("%d.%d.%d.1", overlayFirst, host4[2], host4[3]))
	return gw, nil
}

// ComputeBridgeName returns the fan bridge name for a given overlay network.
//
// Example: overlayNetwork="240.0.0.0/8" → "fan-240"
func ComputeBridgeName(overlayNetwork string) (string, error) {
	overlayFirst, err := validateOverlay(overlayNetwork)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("fan-%d", overlayFirst), nil
}

// ComputeUnderlayArg returns the string passed to fanctl -u, e.g. "172.16.3.4/16".
func ComputeUnderlayArg(hostIP net.IP, underlayPrefix int) string {
	return fmt.Sprintf("%s/%d", hostIP.String(), underlayPrefix)
}

// validateOverlay parses the overlay CIDR and returns the first octet.
// It returns an error if the CIDR is invalid or the network is not IPv4.
func validateOverlay(overlayNetwork string) (byte, error) {
	ip, _, err := net.ParseCIDR(overlayNetwork)
	if err != nil {
		return 0, fmt.Errorf("invalid overlay network %q: %w", overlayNetwork, err)
	}

	overlay4 := ip.To4()
	if overlay4 == nil {
		return 0, fmt.Errorf("overlay network must be IPv4, got: %s", overlayNetwork)
	}

	return overlay4[0], nil
}

// to4 converts a net.IP to its 4-byte IPv4 form, or returns nil if it is not IPv4.
func to4(ip net.IP) net.IP {
	if ip == nil {
		return nil
	}
	return ip.To4()
}
