// Package config provides CNI configuration parsing for the fancni plugin.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

const (
	defaultOverlayNetwork = "240.0.0.0/8"
	defaultUnderlayPrefix = 16
)

// NetConfig holds the CNI configuration fields for fancni.
type NetConfig struct {
	CNIVersion     string `json:"cniVersion"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	OverlayNetwork string `json:"overlayNetwork"`
	UnderlayPrefix int    `json:"underlayPrefix"`
}

// Parse reads a CNI JSON config from r, applies defaults, and validates the
// resulting configuration.
//
// Defaults applied:
//   - overlayNetwork defaults to "240.0.0.0/8" when empty
//   - underlayPrefix defaults to 16 when zero
func Parse(r io.Reader) (NetConfig, error) {
	var cfg NetConfig

	data, err := io.ReadAll(r)
	if err != nil {
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config JSON: %w", err)
	}

	// Apply defaults.
	if cfg.OverlayNetwork == "" {
		cfg.OverlayNetwork = defaultOverlayNetwork
	}
	if cfg.UnderlayPrefix == 0 {
		cfg.UnderlayPrefix = defaultUnderlayPrefix
	}

	// Validate.
	if err := validate(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// validate checks that the parsed (and defaulted) config is semantically valid.
func validate(cfg NetConfig) error {
	ip, _, err := net.ParseCIDR(cfg.OverlayNetwork)
	if err != nil {
		return fmt.Errorf("overlayNetwork %q is not a valid CIDR: %w", cfg.OverlayNetwork, err)
	}
	if ip.To4() == nil {
		return fmt.Errorf("overlayNetwork %q must be an IPv4 CIDR", cfg.OverlayNetwork)
	}
	if cfg.UnderlayPrefix < 1 || cfg.UnderlayPrefix > 32 {
		return fmt.Errorf("underlayPrefix must be between 1 and 32, got %d", cfg.UnderlayPrefix)
	}
	return nil
}
