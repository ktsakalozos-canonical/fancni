package netutil_test

import (
	"testing"

	"github.com/ktsakalozos-canonical/fancni/internal/netutil"
)

func TestLinkExists_NonExistent(t *testing.T) {
	// A link with this name should never exist in a normal test environment.
	if netutil.LinkExists("fancni-test-nonexistent-veth9x") {
		t.Error("expected LinkExists to return false for nonexistent interface")
	}
}

func TestCreateVethPair_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestAttachToBridge_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestMoveToNetNS_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestConfigureInterface_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestDeleteLink_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestGetLinkMAC_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}

func TestVerifyInterfaceConfig_RequiresRoot(t *testing.T) {
	t.Skip("requires root")
}
