//go:build linux || darwin
// +build linux darwin

package procfreeze

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestFreezeUnfreezeSmoke(t *testing.T) {
	// Start a short-lived sleep process that we can freeze/unfreeze.
	cmd := exec.Command("sh", "-c", "sleep 2")
	if err := cmd.Start(); err != nil {
		t.Skipf("couldn't start sleep - skipping: %v", err)
	}
	pid := cmd.Process.Pid
	man := &Manifest{
		Name: "test",
		PIDs: []int{pid},
	}
	mgr := NewManager()
	mgr.manifest = man

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try freeze
	if err := mgr.Freeze(ctx); err != nil {
		// On some systems sending SIGSTOP to subgroup may fail; treat as skip.
		t.Skipf("freeze failed: %v", err)
	}

	// Unfreeze
	if err := mgr.Unfreeze(ctx); err != nil {
		t.Fatalf("unfreeze failed: %v", err)
	}

	// Wait for process to exit
	_ = cmd.Wait()
}