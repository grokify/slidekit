package ops

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/slidekit/backends/marp"
	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
)

func init() {
	// Register marp backend for all tests
	DefaultRegistry.Register("marp", marp.NewBackend())
}

func TestReadDeck(t *testing.T) {
	// Create a temp file with a simple presentation
	content := `---
marp: true
theme: default
---

# Test Presentation

By Test Author

---

## Slide Two

- Point one
- Point two
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	// Test TOON format
	result, err := ReadDeck(ctx, ref, ReadOptions{Format: format.FormatTOON})
	if err != nil {
		t.Fatalf("ReadDeck failed: %v", err)
	}

	if result.Deck == nil {
		t.Fatal("expected non-nil deck")
	}
	if result.Deck.Title != "Test Presentation" {
		t.Errorf("expected title 'Test Presentation', got %q", result.Deck.Title)
	}
	if !strings.Contains(result.Output, "deck Test Presentation") {
		t.Errorf("TOON output missing deck header: %s", result.Output)
	}

	// Test JSON format
	result, err = ReadDeck(ctx, ref, ReadOptions{Format: format.FormatJSON})
	if err != nil {
		t.Fatalf("ReadDeck with JSON failed: %v", err)
	}
	if !strings.Contains(result.Output, `"title": "Test Presentation"`) {
		t.Errorf("JSON output missing title: %s", result.Output)
	}
}

func TestReadDeckFromPath(t *testing.T) {
	content := `---
marp: true
---

# Quick Test
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "quick.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	result, err := ReadDeckFromPath(ctx, path, ReadOptions{})
	if err != nil {
		t.Fatalf("ReadDeckFromPath failed: %v", err)
	}

	if result.Deck.Title != "Quick Test" {
		t.Errorf("expected title 'Quick Test', got %q", result.Deck.Title)
	}
}

func TestReadDeckNotFound(t *testing.T) {
	ctx := context.Background()
	_, err := ReadDeckFromPath(ctx, "/nonexistent/file.md", ReadOptions{})
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadDeckWithSections(t *testing.T) {
	content := `---
marp: true
---

# Main Title

---

<!-- _class: section-divider -->

# Section One

---

## Content Slide

Some content here.

---

<!-- _class: section-divider -->

# Section Two

---

## Another Slide

More content.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "sections.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	result, err := ReadDeckFromPath(ctx, path, ReadOptions{})
	if err != nil {
		t.Fatalf("ReadDeckFromPath failed: %v", err)
	}

	// Should have sections
	if len(result.Deck.Sections) < 2 {
		t.Errorf("expected at least 2 sections, got %d", len(result.Deck.Sections))
	}

	// Count total slides
	totalSlides := result.Deck.SlideCount()
	if totalSlides != 5 {
		t.Errorf("expected 5 slides, got %d", totalSlides)
	}
}

func TestReadDeckWithNotes(t *testing.T) {
	content := `---
marp: true
---

# Slide With Notes

Some content.

<!--
These are speaker notes.
They span multiple lines.
-->
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "notes.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	result, err := ReadDeckFromPath(ctx, path, ReadOptions{})
	if err != nil {
		t.Fatalf("ReadDeckFromPath failed: %v", err)
	}

	slides := result.Deck.AllSlides()
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}

	if !slides[0].HasNotes() {
		t.Error("expected slide to have notes")
	}

	notesText := slides[0].NotesText()
	if !strings.Contains(notesText, "speaker notes") {
		t.Errorf("notes missing expected content: %s", notesText)
	}
}
