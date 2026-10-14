package ramguard

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRunWithSecretsBasic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secrets := map[string]string{
		"SUBSTRATE_TEST_SECRET": "s3cr3t-value",
	}

	var path string
	var args []string

	if runtime.GOOS == "windows" {
		path = "cmd.exe"
		args = []string{"/C", "echo", "ok"}
	} else {
		// prefer a stable binary
		if _, err := os.Stat("/usr/bin/env"); err == nil {
			path = "/usr/bin/env"
			args = []string{"true"}
		} else {
			path = "true"
			args = nil
		}
	}

	if err := RunWithSecrets(ctx, path, args, secrets); err != nil {
		t.Fatalf("RunWithSecrets returned error: %v", err)
	}
}