//go:build windows
// +build windows

package procfreeze

import (
	"context"
	"testing"
	"time"
)

func TestWindowsNotSupported(t *testing.T) {
	mgr := NewManager()
	mgr.manifest = &Manifest{
		Name: "wtest",
		PIDs: []int{1},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// On Windows we expect ErrUnsupportedOS if freezePID/unfreezePID have not been implemented.
	if err := mgr.Freeze(ctx); err == nil {
		t.Fatalf("expected unsupported error on windows Freeze")
	}
	if err := mgr.Unfreeze(ctx); err == nil {
		t.Fatalf("expected unsupported error on windows Unfreeze")
	}
}