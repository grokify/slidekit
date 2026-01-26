package marp

import (
	"os"
	"testing"

	"github.com/grokify/slidekit/format"
)

func TestRealPresentationFile(t *testing.T) {
	path := "/Users/johnwang/go/src/github.com/agentplexus/stats-agent-team/PRESENTATION.md"

	// Skip if file doesn't exist (CI environments)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("real presentation file not found, skipping")
	}

	reader := NewReader()
	deck, err := reader.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	// Basic structural checks
	if deck.Title == "" {
		t.Error("expected non-empty deck title")
	}
	t.Logf("Title: %s", deck.Title)

	if deck.Theme == nil {
		t.Error("expected non-nil theme")
	} else {
		t.Logf("Theme: %s", deck.Theme.Name)
	}

	t.Logf("Sections: %d", len(deck.Sections))
	t.Logf("Total slides: %d", deck.SlideCount())

	if len(deck.Sections) < 5 {
		t.Errorf("expected at least 5 sections, got %d", len(deck.Sections))
	}
	if deck.SlideCount() < 30 {
		t.Errorf("expected at least 30 slides, got %d", deck.SlideCount())
	}

	// Check sections have titles
	for _, s := range deck.Sections {
		t.Logf("  Section %q: %d slides", s.Title, len(s.Slides))
	}

	// Check speaker notes are parsed
	notesCount := 0
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if len(slide.Notes) > 0 {
				notesCount++
			}
		}
	}
	t.Logf("Slides with notes: %d", notesCount)
	if notesCount < 20 {
		t.Errorf("expected at least 20 slides with notes, got %d", notesCount)
	}

	// Verify TOON encoding works
	encoder := format.NewTOONEncoder()
	toon := encoder.EncodeDeck(deck)
	if len(toon) < 1000 {
		t.Errorf("expected substantial TOON output, got %d bytes", len(toon))
	}
	t.Logf("TOON output: %d bytes", len(toon))
}
