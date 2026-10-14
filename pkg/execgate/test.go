package execgate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanDependencyManifestsAndForceRun(t *testing.T) {
	// Create a temporary directory with a dummy manifest
	dir := t.TempDir()
	fn := filepath.Join(dir, "package.json")
	if err := os.WriteFile(fn, []byte(`{"name":"test"}`), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Use a harmless command and force=true so we don't require interactive input.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Craft a command that matches risk signature (rm -rf) but actually is harmless (echo)
	cmd := "rm -rf ./does-not-exist || echo ok"
	executed, err := InspectAndMaybeRun(ctx, cmd, true)
	if err != nil {
		t.Fatalf("InspectAndMaybeRun failed: %v", err)
	}
	if !executed {
		t.Fatalf("expected executed=true")
	}
}