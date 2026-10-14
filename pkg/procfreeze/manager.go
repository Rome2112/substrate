package procfreeze

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Manager manages a workspace manifest and provides Freeze/Unfreeze operations.
type Manager struct {
	manifest *Manifest
}

// NewManager constructs a Manager.
func NewManager() *Manager {
	return &Manager{}
}

// LoadManifest reads and parses the manifest JSON file at path.
func (m *Manager) LoadManifest(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open manifest: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var manifest Manifest
	if err := dec.Decode(&manifest); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}
	m.manifest = &manifest
	return nil
}

// Freeze suspends all PIDs listed in the manifest using platform primitives.
func (m *Manager) Freeze(ctx context.Context) error {
	if m.manifest == nil {
		return ErrManifestNotLoaded
	}
	// Provide a simple timeout context for each operation
	for _, pid := range m.manifest.PIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := freezePID(pid); err != nil {
			// Best-effort: continue and aggregate errors
			return fmt.Errorf("freeze pid %d: %w", pid, err)
		}
		// brief delay so large groups don't overload scheduler
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

// Unfreeze resumes all PIDs listed in the manifest using platform primitives.
func (m *Manager) Unfreeze(ctx context.Context) error {
	if m.manifest == nil {
		return ErrManifestNotLoaded
	}
	for _, pid := range m.manifest.PIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := unfreezePID(pid); err != nil {
			return fmt.Errorf("unfreeze pid %d: %w", pid, err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}