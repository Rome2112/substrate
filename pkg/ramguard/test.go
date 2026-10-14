package ramguard

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRunWithSecretsBasic(t *testing.T) {
	// Use cross-platform harmless command.
	cmd := "true"
	if runtime.GOOS == "windows" {
		cmd = "cmd"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secrets := map[string]string{
		"SUBSTRATE_TEST_SECRET": "s3cr3t-value",
	}
	var path string
	var args []string
	if runtime.GOOS == "windows" {
		// Use echo in cmd
		path = "cmd.exe"
		args = []string{"/C", "echo", "ok"}
	} else {
		path = "/usr/bin/true"
		args = []string{}
		if _, err := os.Stat(path); err != nil {
			// fallback to "true" in PATH
			path = "true"
		}
	}

	err := RunWithSecrets(ctx, path, args, secrets)
	if err != nil {
		t.Fatalf("RunWithSecrets returned error: %v", err)
	}
}