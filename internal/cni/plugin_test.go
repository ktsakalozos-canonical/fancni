package cni

import (
	"bytes"
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/ktsakalozos-canonical/fancni/internal/config"
)

// captureStdout replaces os.Stdout with a pipe, runs f, and returns the
// captured output.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	f()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// TestHandleVersion checks that HandleVersion writes valid JSON with cniVersion
// 1.0.0 and supportedVersions containing "1.0.0". No root needed.
func TestHandleVersion(t *testing.T) {
	p := &Plugin{} // HandleVersion needs no fields.

	out := captureStdout(t, func() {
		if err := p.HandleVersion(); err != nil {
			t.Fatalf("HandleVersion: %v", err)
		}
	})

	var v struct {
		CNIVersion        string   `json:"cniVersion"`
		SupportedVersions []string `json:"supportedVersions"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("unmarshal output %q: %v", out, err)
	}
	if v.CNIVersion != "1.0.0" {
		t.Errorf("cniVersion = %q, want 1.0.0", v.CNIVersion)
	}
	if len(v.SupportedVersions) == 0 || v.SupportedVersions[0] != "1.0.0" {
		t.Errorf("supportedVersions = %v, want [1.0.0]", v.SupportedVersions)
	}
}

// TestCNIResultJSONMarshal validates that the cniResult struct marshals to the
// expected CNI 1.0.0 JSON shape (no "version" field inside ips entries).
func TestCNIResultJSONMarshal(t *testing.T) {
	result := cniResult{
		CNIVersion: "1.0.0",
		Interfaces: []cniInterface{
			{Name: "eth0", MAC: "aa:bb:cc:dd:ee:ff", Sandbox: "/var/run/netns/xxx"},
		},
		IPs: []cniIP{
			{Address: "240.3.4.2/24", Gateway: "240.3.4.1", Interface: 0},
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Unmarshal into a generic map to verify structure.
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m["cniVersion"] != "1.0.0" {
		t.Errorf("cniVersion = %v, want 1.0.0", m["cniVersion"])
	}

	ips, ok := m["ips"].([]interface{})
	if !ok || len(ips) == 0 {
		t.Fatalf("ips field missing or empty")
	}
	ipEntry, ok := ips[0].(map[string]interface{})
	if !ok {
		t.Fatalf("ips[0] not an object")
	}
	// CNI 1.0.0 must NOT have a "version" field in ip entries.
	if _, hasVersion := ipEntry["version"]; hasVersion {
		t.Errorf("ips[0] has unexpected 'version' field (CNI 1.0.0 removed it)")
	}
	if ipEntry["address"] != "240.3.4.2/24" {
		t.Errorf("ips[0].address = %v, want 240.3.4.2/24", ipEntry["address"])
	}
	if ipEntry["gateway"] != "240.3.4.1" {
		t.Errorf("ips[0].gateway = %v, want 240.3.4.1", ipEntry["gateway"])
	}
}

// TestRandomVethName verifies that generated veth names have the right prefix
// and length and are random enough.
func TestRandomVethName(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 20; i++ {
		name, err := randomVethName()
		if err != nil {
			t.Fatalf("randomVethName: %v", err)
		}
		if len(name) != 8 { // "veth" (4) + 4 hex chars
			t.Errorf("name %q has length %d, want 8", name, len(name))
		}
		if name[:4] != "veth" {
			t.Errorf("name %q does not start with 'veth'", name)
		}
		seen[name] = true
	}
	// With 2 random bytes there are 65536 possibilities; 20 draws should
	// produce at least 2 distinct values with overwhelming probability.
	if len(seen) < 2 {
		t.Errorf("got only %d distinct veth names in 20 tries — likely not random", len(seen))
	}
}

// TestHandleDelIdempotent checks that DEL with no IPAM entry returns nil.
func TestHandleDelIdempotent(t *testing.T) {
	p := &Plugin{
		config:      config.NetConfig{OverlayNetwork: "240.0.0.0/8", UnderlayPrefix: 16},
		ipam:        &stubIPAM{},
		hostIP:      net.ParseIP("172.16.3.4"),
		containerID: "nonexistent",
	}
	if err := p.HandleDel(); err != nil {
		t.Errorf("HandleDel with missing container: %v", err)
	}
}

// stubIPAM is a minimal IPAM that returns not-found for everything.
type stubIPAM struct{}

func (s *stubIPAM) Allocate(_ string) (net.IP, error)           { return nil, nil }
func (s *stubIPAM) Lookup(_ string) (net.IP, bool, error)       { return nil, false, nil }
func (s *stubIPAM) Free(_ string) (net.IP, bool, error)         { return nil, false, nil }
