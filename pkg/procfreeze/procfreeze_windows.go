//go:build windows
// +build windows

package procfreeze

import "fmt"

// Windows stub implementation — returns ErrUnsupportedOS until a full Job Object/thread
// suspension implementation is provided.
func freezePID(pid int) error {
	return fmt.Errorf("freezePID not implemented on Windows: %w", ErrUnsupportedOS)
}

func unfreezePID(pid int) error {
	return fmt.Errorf("unfreezePID not implemented on Windows: %w", ErrUnsupportedOS)
}

// AttachPIDsToJobObject is a placeholder for potential Windows Job Object support.
func AttachPIDsToJobObject(jobName string, pids []int) error {
	_ = jobName
	_ = pids
	return fmt.Errorf("Windows Job Object support not implemented: %w", ErrUnsupportedOS)
}