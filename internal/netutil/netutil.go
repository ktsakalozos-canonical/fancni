// Package netutil provides netlink-based networking helpers for fancni.
// All operations use the netlink library; no exec calls to "ip" are made.
package netutil

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// CreateVethPair creates a veth pair with the given names.
// If a veth with hostVethName already exists the call is idempotent.
func CreateVethPair(hostVethName, containerVethName string) error {
	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: hostVethName},
		PeerName:  containerVethName,
	}
	if err := netlink.LinkAdd(veth); err != nil {
		// EEXIST — already created, treat as success.
		if isExist(err) {
			return nil
		}
		return fmt.Errorf("creating veth pair %s/%s: %w", hostVethName, containerVethName, err)
	}
	return nil
}

// AttachToBridge sets ifName's master to bridgeName and brings ifName up.
func AttachToBridge(ifName, bridgeName string) error {
	bridge, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("getting bridge %q: %w", bridgeName, err)
	}
	iface, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("getting interface %q: %w", ifName, err)
	}
	if err := netlink.LinkSetMaster(iface, bridge); err != nil {
		return fmt.Errorf("setting master of %q to %q: %w", ifName, bridgeName, err)
	}
	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("bringing up %q: %w", ifName, err)
	}
	return nil
}

// MoveToNetNS moves ifName into the namespace identified by nsFd.
func MoveToNetNS(ifName string, nsFd int) error {
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("getting link %q: %w", ifName, err)
	}
	if err := netlink.LinkSetNsFd(link, nsFd); err != nil {
		return fmt.Errorf("moving %q to netns fd %d: %w", ifName, nsFd, err)
	}
	return nil
}

// ConfigureInterface opens nsPath, switches into it, configures ifName with
// ipCIDR, and adds a default route via gatewayIP.
// The original namespace is always restored before returning.
func ConfigureInterface(nsPath, ifName, ipCIDR, gatewayIP string) error {
	// Lock OS thread so namespace operations are thread-safe.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origNS, err := netns.Get()
	if err != nil {
		return fmt.Errorf("getting current netns: %w", err)
	}
	defer origNS.Close()

	targetNS, err := netns.GetFromPath(nsPath)
	if err != nil {
		return fmt.Errorf("opening netns %q: %w", nsPath, err)
	}
	defer targetNS.Close()

	if err := netns.Set(targetNS); err != nil {
		return fmt.Errorf("switching to netns %q: %w", nsPath, err)
	}
	defer func() { _ = netns.Set(origNS) }()

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("getting link %q in ns %q: %w", ifName, nsPath, err)
	}

	mac, err := generateMAC()
	if err != nil {
		return err
	}
	if err := netlink.LinkSetHardwareAddr(link, mac); err != nil {
		return fmt.Errorf("setting MAC on %q: %w", ifName, err)
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("bringing up %q: %w", ifName, err)
	}

	addr, err := netlink.ParseAddr(ipCIDR)
	if err != nil {
		return fmt.Errorf("parsing address %q: %w", ipCIDR, err)
	}
	if err := netlink.AddrAdd(link, addr); err != nil && !isExist(err) {
		return fmt.Errorf("adding address %q to %q: %w", ipCIDR, ifName, err)
	}

	gw := net.ParseIP(gatewayIP)
	if gw == nil {
		return fmt.Errorf("invalid gateway IP %q", gatewayIP)
	}
	// Default route: 0.0.0.0/0
	_, dst, _ := net.ParseCIDR("0.0.0.0/0")
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Gw:        gw,
	}
	if err := netlink.RouteAdd(route); err != nil && !isExist(err) {
		return fmt.Errorf("adding default route via %q: %w", gatewayIP, err)
	}

	return nil
}

// DeleteLink deletes the named link. Returns nil if the link does not exist.
func DeleteLink(ifName string) error {
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		if isNotFound(err) {
			return nil
		}
		return fmt.Errorf("getting link %q: %w", ifName, err)
	}
	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("deleting link %q: %w", ifName, err)
	}
	return nil
}

// GetLinkMAC opens nsPath and returns the hardware address of ifName.
func GetLinkMAC(nsPath, ifName string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origNS, err := netns.Get()
	if err != nil {
		return "", fmt.Errorf("getting current netns: %w", err)
	}
	defer origNS.Close()

	targetNS, err := netns.GetFromPath(nsPath)
	if err != nil {
		return "", fmt.Errorf("opening netns %q: %w", nsPath, err)
	}
	defer targetNS.Close()

	if err := netns.Set(targetNS); err != nil {
		return "", fmt.Errorf("switching to netns %q: %w", nsPath, err)
	}
	defer func() { _ = netns.Set(origNS) }()

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return "", fmt.Errorf("getting link %q: %w", ifName, err)
	}
	return link.Attrs().HardwareAddr.String(), nil
}

// LinkExists reports whether a link with the given name exists in the current
// network namespace.
func LinkExists(ifName string) bool {
	_, err := netlink.LinkByName(ifName)
	return err == nil
}

// VerifyInterfaceConfig opens nsPath and checks that ifName is up and has
// expectedIPCIDR assigned. Returns a descriptive error or nil.
func VerifyInterfaceConfig(nsPath, ifName, expectedIPCIDR string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origNS, err := netns.Get()
	if err != nil {
		return fmt.Errorf("getting current netns: %w", err)
	}
	defer origNS.Close()

	targetNS, err := netns.GetFromPath(nsPath)
	if err != nil {
		return fmt.Errorf("opening netns %q: %w", nsPath, err)
	}
	defer targetNS.Close()

	if err := netns.Set(targetNS); err != nil {
		return fmt.Errorf("switching to netns %q: %w", nsPath, err)
	}
	defer func() { _ = netns.Set(origNS) }()

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("interface %q not found: %w", ifName, err)
	}
	if link.Attrs().Flags&net.FlagUp == 0 {
		return fmt.Errorf("interface %q is not up", ifName)
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return fmt.Errorf("listing addresses on %q: %w", ifName, err)
	}
	for _, a := range addrs {
		if a.IPNet != nil && a.IPNet.String() == expectedIPCIDR {
			return nil
		}
	}
	return fmt.Errorf("interface %q does not have address %q", ifName, expectedIPCIDR)
}

// generateMAC returns a random locally-administered unicast MAC address.
func generateMAC() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("generating random MAC: %w", err)
	}
	// Set locally-administered bit, clear multicast bit
	buf[0] = (buf[0] | 0x02) & 0xfe
	return net.HardwareAddr(buf), nil
}

// isExist returns true for "already exists" errors (EEXIST).
func isExist(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, os.ErrExist)
}

// isNotFound returns true for link-not-found errors.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(netlink.LinkNotFoundError)
	return ok
}
