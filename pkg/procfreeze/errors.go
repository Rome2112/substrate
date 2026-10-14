package procfreeze

import "errors"

var (
	// ErrUnsupportedOS is returned when a requested operation isn't implemented on the current OS.
	ErrUnsupportedOS = errors.New("operation not supported on this OS")
	// ErrManifestNotLoaded indicates operations were attempted before a manifest was loaded.
	ErrManifestNotLoaded = errors.New("manifest not loaded")
)