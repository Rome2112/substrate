package execgate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestContainsRisk(t *testing.T) {
	cases := []struct {
		cmd  string
		want bool
	}{
		{"echo hello", false},
		{"rm -rf ./dist", true},
		{"psql postgres://prod-db/app", true},
		{"go test ./...", false},
	}

	for _, tc := range cases {
		if got := containsRisk(tc.cmd); got != tc.want {
			t.Fatalf("containsRisk(%q) = %v, want %v", tc.cmd, got, tc.want)
		}
	}
}

func TestScanDependencyManifests(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"go.mod", "package.json", ".git"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.Dir(name)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, name), []byte("test"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	paths := scanDependencyManifests(root)
	if len(paths) == 0 {
		t.Fatal("expected dependency manifests to be found")
	}
}

func TestInspectAndMaybeRunForce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ok, err := InspectAndMaybeRun(ctx, "echo substrate ok", true)
	if err != nil {
		t.Fatalf("InspectAndMaybeRun returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected execution to continue with force=true")
	}
}