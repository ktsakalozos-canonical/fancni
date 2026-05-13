// Package cni implements the CNI plugin orchestration layer for fancni.
package cni

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/ktsakalozos-canonical/fancni/internal/config"
	"github.com/ktsakalozos-canonical/fancni/internal/fan"
	"github.com/ktsakalozos-canonical/fancni/internal/ipam"
	"github.com/ktsakalozos-canonical/fancni/internal/netutil"
)

// Plugin orchestrates CNI ADD/DEL/CHECK/VERSION operations.
type Plugin struct {
	config   config.NetConfig
	ipam     ipam.IPAM
	hostIP   net.IP

	// From CNI env vars
	command     string
	ifName      string
	netNS       string
	containerID string
}

// NewPlugin creates a new Plugin, reading CNI env vars from the environment.
func NewPlugin(cfg config.NetConfig, i ipam.IPAM, hostIP net.IP) *Plugin {
	return &Plugin{
		config:      cfg,
		ipam:        i,
		hostIP:      hostIP,
		command:     os.Getenv("CNI_COMMAND"),
		ifName:      os.Getenv("CNI_IFNAME"),
		netNS:       os.Getenv("CNI_NETNS"),
		containerID: os.Getenv("CNI_CONTAINERID"),
	}
}

// HandleAdd performs the CNI ADD operation.
func (p *Plugin) HandleAdd() error {
	if p.netNS == "" {
		return fmt.Errorf("CNI_NETNS is required for ADD")
	}

	// 1. Compute fan gateway and bridge name.
	gw, err := fan.ComputeGateway(p.config.OverlayNetwork, p.hostIP)
	if err != nil {
		return fmt.Errorf("computing gateway: %w", err)
	}
	bridgeName, err := fan.ComputeBridgeName(p.config.OverlayNetwork)
	if err != nil {
		return fmt.Errorf("computing bridge name: %w", err)
	}

	// 2. Ensure fan bridge exists.
	if err := fan.EnsureBridge(bridgeName, p.config.OverlayNetwork, p.hostIP, p.config.UnderlayPrefix); err != nil {
		return fmt.Errorf("ensuring fan bridge: %w", err)
	}

	// 3. Allocate IP via IPAM.
	ip, err := p.ipam.Allocate(p.containerID)
	if err != nil {
		return fmt.Errorf("allocating IP: %w", err)
	}

	// Build IP CIDR for the container (/24 matches the fan pod subnet).
	ipCIDR := fmt.Sprintf("%s/24", ip.String())

	// 4. Generate veth names.
	hostVethName, err := randomVethName()
	if err != nil {
		return fmt.Errorf("generating veth name: %w", err)
	}
	containerVethName := p.ifName

	// 5. Create veth pair.
	if err := netutil.CreateVethPair(hostVethName, containerVethName); err != nil {
		return fmt.Errorf("creating veth pair: %w", err)
	}

	// 6. Attach host veth to fan bridge.
	if err := netutil.AttachToBridge(hostVethName, bridgeName); err != nil {
		return fmt.Errorf("attaching veth to bridge: %w", err)
	}

	// 7. Move container veth to netns.
	nsFile, err := os.Open(p.netNS)
	if err != nil {
		return fmt.Errorf("opening netns %q: %w", p.netNS, err)
	}
	defer nsFile.Close()

	if err := netutil.MoveToNetNS(containerVethName, int(nsFile.Fd())); err != nil {
		return fmt.Errorf("moving veth to netns: %w", err)
	}

	// 8. Configure interface inside netns.
	if err := netutil.ConfigureInterface(p.netNS, containerVethName, ipCIDR, gw.String()); err != nil {
		return fmt.Errorf("configuring interface: %w", err)
	}

	// 9. Get MAC address.
	mac, err := netutil.GetLinkMAC(p.netNS, containerVethName)
	if err != nil {
		return fmt.Errorf("getting MAC: %w", err)
	}

	// 10. Build and print CNI result JSON.
	result := cniResult{
		CNIVersion: "1.0.0",
		Interfaces: []cniInterface{
			{
				Name:    containerVethName,
				MAC:     mac,
				Sandbox: p.netNS,
			},
		},
		IPs: []cniIP{
			{
				Address:   ipCIDR,
				Gateway:   gw.String(),
				Interface: 0,
			},
		},
	}
	return writeJSON(result)
}

// HandleDel performs the CNI DEL operation.
func (p *Plugin) HandleDel() error {
	// 1. Lookup container IP in IPAM.
	_, found, err := p.ipam.Lookup(p.containerID)
	if err != nil {
		return fmt.Errorf("IPAM lookup: %w", err)
	}
	// 2. If not found, return nil (idempotent).
	if !found {
		return nil
	}

	// 3. Free IP in IPAM. Veth cleanup is automatic when the netns is destroyed.
	if _, _, err := p.ipam.Free(p.containerID); err != nil {
		return fmt.Errorf("IPAM free: %w", err)
	}
	return nil
}

// HandleCheck performs the CNI CHECK operation.
func (p *Plugin) HandleCheck() error {
	if p.netNS == "" {
		return fmt.Errorf("CNI_NETNS is required for CHECK")
	}

	// 1. Lookup container IP in IPAM.
	ip, found, err := p.ipam.Lookup(p.containerID)
	if err != nil {
		return fmt.Errorf("IPAM lookup: %w", err)
	}
	if !found {
		return fmt.Errorf("container %q has no IPAM allocation", p.containerID)
	}

	// 2. Compute expected fan bridge name (verify bridge exists).
	bridgeName, err := fan.ComputeBridgeName(p.config.OverlayNetwork)
	if err != nil {
		return fmt.Errorf("computing bridge name: %w", err)
	}
	if !netutil.LinkExists(bridgeName) {
		return fmt.Errorf("fan bridge %q does not exist", bridgeName)
	}

	// 3. Verify interface config inside container netns.
	ipCIDR := fmt.Sprintf("%s/24", ip.String())
	if err := netutil.VerifyInterfaceConfig(p.netNS, p.ifName, ipCIDR); err != nil {
		return fmt.Errorf("interface check failed: %w", err)
	}
	return nil
}

// HandleVersion prints the CNI version info to stdout.
func (p *Plugin) HandleVersion() error {
	v := struct {
		CNIVersion        string   `json:"cniVersion"`
		SupportedVersions []string `json:"supportedVersions"`
	}{
		CNIVersion:        "1.0.0",
		SupportedVersions: []string{"1.0.0"},
	}
	return writeJSON(v)
}

// --- helpers ---

// cniResult represents the CNI 1.0.0 result JSON.
type cniResult struct {
	CNIVersion string         `json:"cniVersion"`
	Interfaces []cniInterface `json:"interfaces"`
	IPs        []cniIP        `json:"ips"`
}

type cniInterface struct {
	Name    string `json:"name"`
	MAC     string `json:"mac"`
	Sandbox string `json:"sandbox"`
}

type cniIP struct {
	Address   string `json:"address"`
	Gateway   string `json:"gateway"`
	Interface int    `json:"interface"`
}

// writeJSON marshals v and writes it to stdout.
func writeJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	_, err = os.Stdout.Write(data)
	return err
}

// randomVethName generates a veth name like "vethXXXX" using 2 random bytes
// formatted as 4 hex characters.
func randomVethName() (string, error) {
	b := make([]byte, 2)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return fmt.Sprintf("veth%04x", b), nil
}
