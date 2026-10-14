//go:build windows
// +build windows

package procfreeze

import (
	"context"
	"testing"
)

func TestWindowsNotSupported(t *testing.T) {
	mgr := NewManager()
	mgr.manifest = &Manifest{
		Name: "wtest",
		PIDs: []int{1},
	}
	if err := mgr.Freeze(context.Background()); err == nil {
		t.Fatalf("expected unsupported error on windows Freezing")
	}
}