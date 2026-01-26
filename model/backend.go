package model

import "context"

// Backend defines the interface for presentation backends.
type Backend interface {
	// Info returns backend metadata.
	Info() BackendInfo

	// Read loads a presentation from the backend.
	Read(ctx context.Context, ref Ref) (*Deck, error)

	// Plan computes changes needed to reach desired state.
	Plan(ctx context.Context, ref Ref, desired *Deck) (*Diff, error)

	// Apply executes a diff against the backend.
	Apply(ctx context.Context, ref Ref, diff *Diff) error

	// Create creates a new presentation.
	Create(ctx context.Context, deck *Deck) (Ref, error)
}

// Ref identifies a presentation in a backend.
type Ref struct {
	Backend string `json:"backend"` // "marp", "gslides", "reveal"
	ID      string `json:"id"`      // Backend-specific identifier
	Path    string `json:"path"`    // File path (for file-based backends)
}

// BackendInfo describes backend capabilities.
type BackendInfo struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

// Standard capability strings.
const (
	CapabilityRead        = "read"
	CapabilityWrite       = "write"
	CapabilityPlan        = "plan"
	CapabilityApply       = "apply"
	CapabilityCreate      = "create"
	CapabilityTransitions = "transitions"
	CapabilityAudio       = "audio"
	CapabilitySections    = "sections"
)

// HasCapability returns true if the backend has the specified capability.
func (b BackendInfo) HasCapability(cap string) bool {
	for _, c := range b.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}

// IsFileRef returns true if this is a file-based reference.
func (r Ref) IsFileRef() bool {
	return r.Path != ""
}

// String returns a string representation of the reference.
func (r Ref) String() string {
	if r.Path != "" {
		return r.Backend + ":" + r.Path
	}
	return r.Backend + ":" + r.ID
}
