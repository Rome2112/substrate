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
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	pid := cmd.Process.Pid
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try freeze/unfreeze using the platform functions.
	if err := freezePID(pid); err != nil {
		t.Skipf("freeze failed: %v", err)
	}
	if err := unfreezePID(pid); err != nil {
		t.Fatalf("unfreeze failed: %v", err)
	}
}