package ipam

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"syscall"
)

// FileIPAM is a file-backed IPAM implementation that stores container-to-IP
// mappings in a JSON file, guarded by an exclusive flock for concurrency safety.
type FileIPAM struct {
	dataDir  string
	podCIDR  string
	stateFile string
	lockFile  string
}

// NewFileIPAM creates a new FileIPAM.
//
//   - dataDir is the directory where state files are stored (e.g. /var/lib/cni/fancni/).
//   - podCIDR is the node's pod subnet (e.g. 240.3.4.0/24).
//
// The directory is created with MkdirAll if it does not yet exist.
func NewFileIPAM(dataDir, podCIDR string) *FileIPAM {
	return &FileIPAM{
		dataDir:   dataDir,
		podCIDR:   podCIDR,
		stateFile: filepath.Join(dataDir, "ipam.json"),
		lockFile:  filepath.Join(dataDir, "ipam.lock"),
	}
}

// Allocate assigns the next free IP in podCIDR to containerID.
// If containerID already has an allocation the existing IP is returned (idempotent).
func (f *FileIPAM) Allocate(containerID string) (net.IP, error) {
	unlock, err := f.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	allocs, err := f.read()
	if err != nil {
		return nil, err
	}

	// Idempotent: return existing allocation.
	if ipStr, ok := allocs[containerID]; ok {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("ipam: corrupt entry for %s: %q", containerID, ipStr)
		}
		return ip.To4(), nil
	}

	// Build set of already-allocated IPs.
	used := make(map[string]struct{}, len(allocs))
	for _, v := range allocs {
		used[v] = struct{}{}
	}

	// Find the first free IP from .2 to .254.
	ip, err := f.nextFree(used)
	if err != nil {
		return nil, err
	}

	allocs[containerID] = ip.String()
	if err := f.write(allocs); err != nil {
		return nil, err
	}
	return ip, nil
}

// Lookup returns the IP allocated to containerID, or (nil, false, nil) if none.
func (f *FileIPAM) Lookup(containerID string) (net.IP, bool, error) {
	unlock, err := f.lock()
	if err != nil {
		return nil, false, err
	}
	defer unlock()

	allocs, err := f.read()
	if err != nil {
		return nil, false, err
	}

	ipStr, ok := allocs[containerID]
	if !ok {
		return nil, false, nil
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, false, fmt.Errorf("ipam: corrupt entry for %s: %q", containerID, ipStr)
	}
	return ip.To4(), true, nil
}

// Free releases the IP assigned to containerID.
// Returns (nil, false, nil) if containerID has no allocation.
func (f *FileIPAM) Free(containerID string) (net.IP, bool, error) {
	unlock, err := f.lock()
	if err != nil {
		return nil, false, err
	}
	defer unlock()

	allocs, err := f.read()
	if err != nil {
		return nil, false, err
	}

	ipStr, ok := allocs[containerID]
	if !ok {
		return nil, false, nil
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, false, fmt.Errorf("ipam: corrupt entry for %s: %q", containerID, ipStr)
	}

	delete(allocs, containerID)
	if err := f.write(allocs); err != nil {
		return nil, false, err
	}
	return ip.To4(), true, nil
}

// --- internal helpers ---

// lock acquires an exclusive flock on the lock file and returns an unlock func.
func (f *FileIPAM) lock() (func(), error) {
	if err := os.MkdirAll(f.dataDir, 0755); err != nil {
		return nil, fmt.Errorf("ipam: create dataDir: %w", err)
	}
	lf, err := os.OpenFile(f.lockFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("ipam: open lock file: %w", err)
	}
	if err := syscall.Flock(int(lf.Fd()), syscall.LOCK_EX); err != nil {
		lf.Close()
		return nil, fmt.Errorf("ipam: flock: %w", err)
	}
	return func() {
		_ = syscall.Flock(int(lf.Fd()), syscall.LOCK_UN)
		_ = lf.Close()
	}, nil
}

// read loads the allocation map from disk. Returns an empty map if the file
// does not yet exist.
func (f *FileIPAM) read() (map[string]string, error) {
	data, err := os.ReadFile(f.stateFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]string), nil
		}
		return nil, fmt.Errorf("ipam: read state: %w", err)
	}
	var allocs map[string]string
	if err := json.Unmarshal(data, &allocs); err != nil {
		return nil, fmt.Errorf("ipam: parse state: %w", err)
	}
	return allocs, nil
}

// write persists the allocation map atomically (temp file + rename).
func (f *FileIPAM) write(allocs map[string]string) (err error) {
	data, err := json.MarshalIndent(allocs, "", "  ")
	if err != nil {
		return fmt.Errorf("ipam: marshal state: %w", err)
	}

	tmp, err := os.CreateTemp(f.dataDir, "ipam.*.tmp")
	if err != nil {
		return fmt.Errorf("ipam: create temp file: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		return fmt.Errorf("ipam: write temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("ipam: close temp file: %w", err)
	}
	if err = os.Rename(tmp.Name(), f.stateFile); err != nil {
		return fmt.Errorf("ipam: rename temp file: %w", err)
	}
	return nil
}

// nextFree finds the first unallocated IP from .2 to .254 in podCIDR.
func (f *FileIPAM) nextFree(used map[string]struct{}) (net.IP, error) {
	_, ipNet, err := net.ParseCIDR(f.podCIDR)
	if err != nil {
		return nil, fmt.Errorf("ipam: parse podCIDR %q: %w", f.podCIDR, err)
	}

	base := ipNet.IP.To4()
	if base == nil {
		return nil, fmt.Errorf("ipam: podCIDR must be IPv4, got %q", f.podCIDR)
	}

	for i := 2; i <= 254; i++ {
		candidate := net.IP{base[0], base[1], base[2], byte(i)}
		if _, ok := used[candidate.String()]; !ok {
			return candidate, nil
		}
	}
	return nil, fmt.Errorf("ipam: address space exhausted in %s", f.podCIDR)
}
