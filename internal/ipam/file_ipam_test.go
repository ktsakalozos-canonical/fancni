package ipam

import (
	"fmt"
	"net"
	"testing"
)

const testCIDR = "240.3.4.0/24"

func newTestIPAM(t *testing.T) *FileIPAM {
	t.Helper()
	return NewFileIPAM(t.TempDir(), testCIDR)
}

func TestAllocateFirstIP(t *testing.T) {
	f := newTestIPAM(t)
	ip, err := f.Allocate("ctr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := net.IP{240, 3, 4, 2}
	if !ip.Equal(want) {
		t.Errorf("got %v, want %v", ip, want)
	}
}

func TestAllocateSecondIP(t *testing.T) {
	f := newTestIPAM(t)
	if _, err := f.Allocate("ctr-1"); err != nil {
		t.Fatalf("first allocate: %v", err)
	}
	ip, err := f.Allocate("ctr-2")
	if err != nil {
		t.Fatalf("second allocate: %v", err)
	}
	want := net.IP{240, 3, 4, 3}
	if !ip.Equal(want) {
		t.Errorf("got %v, want %v", ip, want)
	}
}

func TestAllocateIdempotent(t *testing.T) {
	f := newTestIPAM(t)
	ip1, err := f.Allocate("ctr-1")
	if err != nil {
		t.Fatalf("first allocate: %v", err)
	}
	ip2, err := f.Allocate("ctr-1")
	if err != nil {
		t.Fatalf("second allocate: %v", err)
	}
	if !ip1.Equal(ip2) {
		t.Errorf("idempotent allocate: got %v, want %v", ip2, ip1)
	}
}

func TestLookupExisting(t *testing.T) {
	f := newTestIPAM(t)
	allocated, _ := f.Allocate("ctr-1")

	ip, found, err := f.Lookup("ctr-1")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !found {
		t.Fatal("expected found=true")
	}
	if !ip.Equal(allocated) {
		t.Errorf("got %v, want %v", ip, allocated)
	}
}

func TestLookupNonExisting(t *testing.T) {
	f := newTestIPAM(t)
	ip, found, err := f.Lookup("no-such-container")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Errorf("expected found=false, got ip=%v", ip)
	}
	if ip != nil {
		t.Errorf("expected nil IP, got %v", ip)
	}
}

func TestFreeExisting(t *testing.T) {
	f := newTestIPAM(t)
	allocated, _ := f.Allocate("ctr-1")

	freed, ok, err := f.Free("ctr-1")
	if err != nil {
		t.Fatalf("free: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if !freed.Equal(allocated) {
		t.Errorf("freed IP %v, want %v", freed, allocated)
	}

	// Subsequent lookup must return not-found.
	_, found, err := f.Lookup("ctr-1")
	if err != nil {
		t.Fatalf("lookup after free: %v", err)
	}
	if found {
		t.Error("expected found=false after free")
	}
}

func TestFreeNonExisting(t *testing.T) {
	f := newTestIPAM(t)
	ip, ok, err := f.Free("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("expected ok=false, got ip=%v", ip)
	}
	if ip != nil {
		t.Errorf("expected nil IP, got %v", ip)
	}
}

func TestAllocateAfterFree(t *testing.T) {
	f := newTestIPAM(t)
	first, _ := f.Allocate("ctr-1")
	f.Free("ctr-1") //nolint:errcheck

	// After freeing the only allocation the next alloc should reuse .2.
	second, err := f.Allocate("ctr-2")
	if err != nil {
		t.Fatalf("allocate after free: %v", err)
	}
	if !second.Equal(first) {
		t.Errorf("expected %v, got %v", first, second)
	}
}

func TestExhaustAddressSpace(t *testing.T) {
	f := newTestIPAM(t)

	// Allocate all 253 usable IPs (.2 – .254).
	for i := 0; i < 253; i++ {
		id := fmt.Sprintf("ctr-%d", i)
		if _, err := f.Allocate(id); err != nil {
			t.Fatalf("allocation %d failed unexpectedly: %v", i, err)
		}
	}

	// 254th allocation must fail.
	_, err := f.Allocate("ctr-overflow")
	if err == nil {
		t.Fatal("expected error when address space exhausted, got nil")
	}
}
