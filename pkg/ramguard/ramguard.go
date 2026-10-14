// Package ramguard provides a minimal API to run a child process with ephemeral,
// in-memory injected environment variables and best-effort zeroing of owned buffers
// after process exit.
//
// Note: This package avoids writing plaintext secrets to disk. It attempts to
// minimize secret lifetime by zeroing the memory buffers it owns after the child
// process has completed. It cannot erase copies made by the OS or other processes.
package ramguard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// ErrMissingPath is returned when an empty path is passed to RunWithSecrets.
var ErrMissingPath = errors.New("executable path is required")

// RunWithSecrets runs the binary at path with args and injects the provided secrets
// as environment variables into the child process's environment block without
// persisting them to disk. On completion it attempts to zero owned buffers.
//
// ctx: cancellation and timeout control.
// path: binary to execute.
// args: arguments (not including path).
// secrets: env map of KEY->VALUE to inject.
func RunWithSecrets(ctx context.Context, path string, args []string, secrets map[string]string) error {
	if strings.TrimSpace(path) == "" {
		return ErrMissingPath
	}

	// Build environment: start from current environment (but we avoid writing files)
	env := os.Environ()
	ownedBuffers := make([][]byte, 0, len(secrets))

	// Append secret env variables in KEY=VALUE form. We capture the byte slices
	// we allocate here in order to zero them after Wait().
	for k, v := range secrets {
		kv := fmt.Sprintf("%s=%s", k, v)
		b := []byte(kv) // owned buffer
		env = append(env, string(b))
		ownedBuffers = append(ownedBuffers, b)
	}

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set low-level process attributes where applicable
	// On Unix we set Setpgid so child runs in its own process group.
	// On Windows syscall.SysProcAttr exists with different fields; we set nothing specialized there.
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	} else {
		// On Windows we prefer to leave defaults (Job object handling is not added).
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	if err := cmd.Start(); err != nil {
		// Attempt to zero owned buffers even on Start failure
		zeroBuffers(ownedBuffers)
		return fmt.Errorf("start failed: %w", err)
	}

	// Wait and then zero.
	err := cmd.Wait()
	zeroBuffers(ownedBuffers)

	if err != nil {
		return fmt.Errorf("process failed: %w", err)
	}
	return nil
}

// zeroBuffers overwrites byte slices with zeros.
func zeroBuffers(buffers [][]byte) {
	for _, b := range buffers {
		for i := range b {
			b[i] = 0
		}
	}
}