package fan

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/ktsakalozos-canonical/fancni/internal/netutil"
)

// EnsureBridge ensures the fan bridge for the given overlay network exists.
// If the bridge already exists, it is a no-op. Otherwise it runs:
//
//	fanctl up -o <overlayNetwork> -u <hostIP>/<underlayPrefix>
//
// This is the only exec.Command call in the entire codebase.
func EnsureBridge(bridgeName, overlayNetwork string, hostIP net.IP, underlayPrefix int) error {
	if netutil.LinkExists(bridgeName) {
		return nil
	}

	underlayArg := ComputeUnderlayArg(hostIP, underlayPrefix)
	cmd := exec.Command("fanctl", "up", "-o", overlayNetwork, "-u", underlayArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Provide a clear error when fanctl is not in PATH.
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return fmt.Errorf("fanctl not found in PATH: install ubuntu-fan package")
		}
		return fmt.Errorf("fanctl up failed: %w\noutput: %s", err, string(out))
	}
	return nil
}
