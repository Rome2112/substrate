package procfreeze

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Manager groups process-level manipulation for a single project context.
type Manager struct {
	manifest *Manifest
}

func NewManager() *Manager {
	return &Manager{}
}

// LoadManifest reads a manifest file from disk and validates it.
func (m *Manager) LoadManifest(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}
	if len(manifest.PIDs) == 0 {
		return ErrNoPIDs
	}
	m.manifest = &manifest
	return nil
}

// Freeze applies suspend operations to each PID in the manifest.
func (m *Manager) Freeze(ctx context.Context) error {
	if m.manifest == nil {
		return ErrManifestNotLoaded
	}
	for _, pid := range m.manifest.PIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := freezePID(pid); err != nil {
			return fmt.Errorf("freeze pid %d: %w", pid, err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	return nil
}

// Unfreeze applies resume operations to each PID in the manifest.
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

func (m *Manager) SetCgroupPath(path string) {
	if m.manifest == nil {
		m.manifest = &Manifest{}
	}
	m.manifest.CgroupPath = path
}
