//go:build linux || darwin
// +build linux darwin

package procfreeze

import (
	"syscall"
)

// freezePID sends SIGSTOP to the given pid (Unix).
func freezePID(pid int) error {
	// Send SIGSTOP to entire process group if possible
	// Use negative pid to address process group
	if err := syscall.Kill(-pid, syscall.SIGSTOP); err != nil {
		// fallback to single pid
		if err2 := syscall.Kill(pid, syscall.SIGSTOP); err2 != nil {
			return err2
		}
	}
	return nil
}

// unfreezePID sends SIGCONT to the given pid (Unix).
func unfreezePID(pid int) error {
	if err := syscall.Kill(-pid, syscall.SIGCONT); err != nil {
		if err2 := syscall.Kill(pid, syscall.SIGCONT); err2 != nil {
			return err2
		}
	}
	return nil
}