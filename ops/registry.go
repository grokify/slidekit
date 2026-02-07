// Package ops provides shared operations for CLI and MCP server.
// Both interfaces use the same underlying logic to ensure consistency.
package ops

import (
	"fmt"
	"sync"

	"github.com/grokify/slidekit/model"
)

// Registry manages backend implementations.
type Registry struct {
	mu       sync.RWMutex
	backends map[string]model.Backend
}

// NewRegistry creates a new backend registry.
func NewRegistry() *Registry {
	return &Registry{
		backends: make(map[string]model.Backend),
	}
}

// Register adds a backend to the registry.
func (r *Registry) Register(name string, backend model.Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends[name] = backend
}

// Get retrieves a backend by name.
func (r *Registry) Get(name string) (model.Backend, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	backend, ok := r.backends[name]
	if !ok {
		return nil, fmt.Errorf("unknown backend: %s", name)
	}
	return backend, nil
}

// List returns all registered backend names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.backends))
	for name := range r.backends {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the global registry used by CLI and MCP.
var DefaultRegistry = NewRegistry()

// DetectBackend determines the backend from a file path.
func DetectBackend(path string) string {
	// For now, default to marp for .md files
	if len(path) > 3 && path[len(path)-3:] == ".md" {
		return "marp"
	}
	return "marp"
}
