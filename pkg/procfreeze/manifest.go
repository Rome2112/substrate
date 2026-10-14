package procfreeze

import "encoding/json"

// Manifest defines a simple workspace grouping, listing PIDs that belong to the
// workspace. It's intentionally small; more metadata can be added later.
type Manifest struct {
	Name string `json:"name"`
	// PIDs is a list of process IDs that belong to this workspace.
	PIDs []int `json:"pids"`
}

// MarshalJSON helper kept for future evolution; currently default works.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	type s Manifest
	return json.Marshal((*s)(m))
}