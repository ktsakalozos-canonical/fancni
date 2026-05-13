// Package iptables provides iptables rule management helpers for fancni.
package iptables

import (
	"fmt"

	goiptables "github.com/coreos/go-iptables/iptables"
)

// EnsureForwardRules adds ACCEPT rules in the FORWARD chain of the filter
// table for traffic to/from podCIDR. Uses AppendUnique so re-running is safe.
func EnsureForwardRules(podCIDR string) error {
	ipt, err := goiptables.New()
	if err != nil {
		return fmt.Errorf("creating iptables handle: %w", err)
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-s", podCIDR, "-j", "ACCEPT"); err != nil {
		return fmt.Errorf("ensuring FORWARD -s %s ACCEPT: %w", podCIDR, err)
	}
	if err := ipt.AppendUnique("filter", "FORWARD", "-d", podCIDR, "-j", "ACCEPT"); err != nil {
		return fmt.Errorf("ensuring FORWARD -d %s ACCEPT: %w", podCIDR, err)
	}
	return nil
}

// EnsureMasqueradeRule adds a MASQUERADE rule in the POSTROUTING chain of the
// nat table for podSubnet traffic not destined for bridgeName.
func EnsureMasqueradeRule(podSubnet, bridgeName string) error {
	ipt, err := goiptables.New()
	if err != nil {
		return fmt.Errorf("creating iptables handle: %w", err)
	}
	if err := ipt.AppendUnique("nat", "POSTROUTING", "-s", podSubnet, "!", "-o", bridgeName, "-j", "MASQUERADE"); err != nil {
		return fmt.Errorf("ensuring MASQUERADE for %s via %s: %w", podSubnet, bridgeName, err)
	}
	return nil
}
