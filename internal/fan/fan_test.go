package fan_test

import (
	"net"
	"testing"

	"github.com/ktsakalozos-canonical/fancni/internal/fan"
)

// parseIP is a test helper that panics on invalid IPs so tests stay concise.
func parseIP(s string) net.IP {
	ip := net.ParseIP(s)
	if ip == nil {
		panic("invalid IP in test: " + s)
	}
	return ip
}

// --- ComputeSubnet ---

func TestComputeSubnet_Standard(t *testing.T) {
	got, err := fan.ComputeSubnet("240.0.0.0/8", parseIP("172.16.3.4"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "240.3.4.0/24" {
		t.Errorf("got %q, want %q", got, "240.3.4.0/24")
	}
}

func TestComputeSubnet_DifferentOverlay(t *testing.T) {
	got, err := fan.ComputeSubnet("10.0.0.0/8", parseIP("172.16.5.6"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "10.5.6.0/24" {
		t.Errorf("got %q, want %q", got, "10.5.6.0/24")
	}
}

func TestComputeSubnet_InvalidCIDR(t *testing.T) {
	_, err := fan.ComputeSubnet("not-a-cidr", parseIP("172.16.3.4"))
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}

func TestComputeSubnet_IPv6Overlay(t *testing.T) {
	_, err := fan.ComputeSubnet("2001:db8::/32", parseIP("172.16.3.4"))
	if err == nil {
		t.Error("expected error for IPv6 overlay, got nil")
	}
}

func TestComputeSubnet_NilHostIP(t *testing.T) {
	_, err := fan.ComputeSubnet("240.0.0.0/8", nil)
	if err == nil {
		t.Error("expected error for nil hostIP, got nil")
	}
}

func TestComputeSubnet_IPv6HostIP(t *testing.T) {
	_, err := fan.ComputeSubnet("240.0.0.0/8", net.ParseIP("2001:db8::1"))
	if err == nil {
		t.Error("expected error for IPv6 hostIP, got nil")
	}
}

// --- ComputeGateway ---

func TestComputeGateway_Standard(t *testing.T) {
	gw, err := fan.ComputeGateway("240.0.0.0/8", parseIP("172.16.3.4"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := net.ParseIP("240.3.4.1")
	if !gw.Equal(want) {
		t.Errorf("got %v, want %v", gw, want)
	}
}

func TestComputeGateway_DifferentOverlay(t *testing.T) {
	gw, err := fan.ComputeGateway("10.0.0.0/8", parseIP("172.16.5.6"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := net.ParseIP("10.5.6.1")
	if !gw.Equal(want) {
		t.Errorf("got %v, want %v", gw, want)
	}
}

func TestComputeGateway_NilHostIP(t *testing.T) {
	_, err := fan.ComputeGateway("240.0.0.0/8", nil)
	if err == nil {
		t.Error("expected error for nil hostIP, got nil")
	}
}

func TestComputeGateway_InvalidCIDR(t *testing.T) {
	_, err := fan.ComputeGateway("bad-cidr", parseIP("172.16.3.4"))
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}

// --- ComputeBridgeName ---

func TestComputeBridgeName_Standard(t *testing.T) {
	got, err := fan.ComputeBridgeName("240.0.0.0/8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fan-240" {
		t.Errorf("got %q, want %q", got, "fan-240")
	}
}

func TestComputeBridgeName_DifferentOverlay(t *testing.T) {
	got, err := fan.ComputeBridgeName("10.0.0.0/8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fan-10" {
		t.Errorf("got %q, want %q", got, "fan-10")
	}
}

func TestComputeBridgeName_InvalidCIDR(t *testing.T) {
	_, err := fan.ComputeBridgeName("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid CIDR, got nil")
	}
}

func TestComputeBridgeName_IPv6(t *testing.T) {
	_, err := fan.ComputeBridgeName("2001:db8::/32")
	if err == nil {
		t.Error("expected error for IPv6 overlay, got nil")
	}
}

// --- ComputeUnderlayArg ---

func TestComputeUnderlayArg(t *testing.T) {
	got := fan.ComputeUnderlayArg(parseIP("172.16.3.4"), 16)
	want := "172.16.3.4/16"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestComputeUnderlayArg_DifferentPrefix(t *testing.T) {
	got := fan.ComputeUnderlayArg(parseIP("10.0.1.2"), 24)
	want := "10.0.1.2/24"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
