package config_test

import (
	"strings"
	"testing"

	"github.com/ktsakalozos-canonical/fancni/internal/config"
)

func TestParse_Valid(t *testing.T) {
	input := `{
		"cniVersion": "1.0.0",
		"name": "fancni",
		"type": "fancni",
		"overlayNetwork": "240.0.0.0/8",
		"underlayPrefix": 16
	}`

	cfg, err := config.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.CNIVersion != "1.0.0" {
		t.Errorf("CNIVersion: got %q, want %q", cfg.CNIVersion, "1.0.0")
	}
	if cfg.Name != "fancni" {
		t.Errorf("Name: got %q, want %q", cfg.Name, "fancni")
	}
	if cfg.OverlayNetwork != "240.0.0.0/8" {
		t.Errorf("OverlayNetwork: got %q, want %q", cfg.OverlayNetwork, "240.0.0.0/8")
	}
	if cfg.UnderlayPrefix != 16 {
		t.Errorf("UnderlayPrefix: got %d, want %d", cfg.UnderlayPrefix, 16)
	}
}

func TestParse_DefaultsApplied(t *testing.T) {
	// overlayNetwork and underlayPrefix omitted — defaults should be applied.
	input := `{"cniVersion": "1.0.0", "name": "fancni", "type": "fancni"}`

	cfg, err := config.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.OverlayNetwork != "240.0.0.0/8" {
		t.Errorf("OverlayNetwork default: got %q, want %q", cfg.OverlayNetwork, "240.0.0.0/8")
	}
	if cfg.UnderlayPrefix != 16 {
		t.Errorf("UnderlayPrefix default: got %d, want %d", cfg.UnderlayPrefix, 16)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := config.Parse(strings.NewReader("{not valid json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParse_InvalidOverlayCIDR(t *testing.T) {
	input := `{"overlayNetwork": "not-a-cidr", "underlayPrefix": 16}`
	_, err := config.Parse(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for invalid overlay CIDR, got nil")
	}
}

func TestParse_IPv6OverlayCIDR(t *testing.T) {
	input := `{"overlayNetwork": "2001:db8::/32", "underlayPrefix": 16}`
	_, err := config.Parse(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for IPv6 overlay CIDR, got nil")
	}
}

func TestParse_NegativeUnderlayPrefix(t *testing.T) {
	input := `{"overlayNetwork": "240.0.0.0/8", "underlayPrefix": -1}`
	_, err := config.Parse(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for negative underlayPrefix, got nil")
	}
}

func TestParse_UnderlayPrefixTooLarge(t *testing.T) {
	input := `{"overlayNetwork": "240.0.0.0/8", "underlayPrefix": 33}`
	_, err := config.Parse(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for underlayPrefix 33 (> 32), got nil")
	}
}

func TestParse_DifferentOverlay(t *testing.T) {
	input := `{"overlayNetwork": "10.0.0.0/8", "underlayPrefix": 24}`
	cfg, err := config.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OverlayNetwork != "10.0.0.0/8" {
		t.Errorf("got %q, want %q", cfg.OverlayNetwork, "10.0.0.0/8")
	}
	if cfg.UnderlayPrefix != 24 {
		t.Errorf("got %d, want %d", cfg.UnderlayPrefix, 24)
	}
}
