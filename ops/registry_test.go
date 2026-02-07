package ops

import (
	"testing"

	"github.com/grokify/slidekit/backends/marp"
)

func TestRegistry(t *testing.T) {
	reg := NewRegistry()

	// Test empty registry
	_, err := reg.Get("marp")
	if err == nil {
		t.Error("expected error for unregistered backend")
	}

	// Register backend
	reg.Register("marp", marp.NewBackend())

	// Test retrieval
	backend, err := reg.Get("marp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend == nil {
		t.Error("expected non-nil backend")
	}

	info := backend.Info()
	if info.Name != "marp" {
		t.Errorf("expected name 'marp', got %q", info.Name)
	}

	// Test List
	names := reg.List()
	if len(names) != 1 {
		t.Errorf("expected 1 backend, got %d", len(names))
	}
	if names[0] != "marp" {
		t.Errorf("expected 'marp', got %q", names[0])
	}
}

func TestDetectBackend(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"presentation.md", "marp"},
		{"slides.md", "marp"},
		{"README.md", "marp"},
		{"file.txt", "marp"}, // defaults to marp
		{"", "marp"},
	}

	for _, tc := range tests {
		got := DetectBackend(tc.path)
		if got != tc.expected {
			t.Errorf("DetectBackend(%q) = %q, want %q", tc.path, got, tc.expected)
		}
	}
}
