package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/slidekit/backends/marp"
	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

func init() {
	// Register marp backend for all tests
	ops.DefaultRegistry.Register("marp", marp.NewBackend())
}

func createTestPresentation(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestHandleReadDeck(t *testing.T) {
	content := `---
marp: true
---

# Test Presentation

Content here.
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	// Test TOON format (default)
	_, output, err := handleReadDeck(ctx, nil, ReadDeckInput{
		Path: path,
	})
	if err != nil {
		t.Fatalf("handleReadDeck failed: %v", err)
	}
	if !strings.Contains(output.Content, "deck Test Presentation") {
		t.Errorf("TOON output missing deck header: %s", output.Content)
	}

	// Test JSON format
	_, output, err = handleReadDeck(ctx, nil, ReadDeckInput{
		Path:   path,
		Format: "json",
	})
	if err != nil {
		t.Fatalf("handleReadDeck with JSON failed: %v", err)
	}
	if !strings.Contains(output.Content, `"title"`) {
		t.Errorf("JSON output missing title: %s", output.Content)
	}
}

func TestHandleReadDeckNotFound(t *testing.T) {
	ctx := context.Background()

	_, _, err := handleReadDeck(ctx, nil, ReadDeckInput{
		Path: "/nonexistent/file.md",
	})
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestHandleListSlides(t *testing.T) {
	content := `---
marp: true
---

# Slide One

---

## Slide Two

---

## Slide Three
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	_, output, err := handleListSlides(ctx, nil, ListSlidesInput{
		Path: path,
	})
	if err != nil {
		t.Fatalf("handleListSlides failed: %v", err)
	}

	// Should list all slides
	if !strings.Contains(output.Content, "slide") {
		t.Errorf("output missing slide info: %s", output.Content)
	}
}

func TestHandleGetSlide(t *testing.T) {
	content := `---
marp: true
---

# First Slide

---

## Target Slide

- Bullet one
- Bullet two
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	// First list to get slide ID
	_, listOutput, err := handleListSlides(ctx, nil, ListSlidesInput{Path: path})
	if err != nil {
		t.Fatalf("handleListSlides failed: %v", err)
	}

	// Get list result to find ID
	listResult, err := ops.ListSlidesFromPath(ctx, path, "toon")
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	if len(listResult.Slides) < 2 {
		t.Fatalf("expected at least 2 slides, got %d", len(listResult.Slides))
	}

	slideID := listResult.Slides[1].ID

	_, output, err := handleGetSlide(ctx, nil, GetSlideInput{
		Path:    path,
		SlideID: slideID,
	})
	if err != nil {
		t.Fatalf("handleGetSlide failed: %v", err)
	}

	if !strings.Contains(output.Content, "Target Slide") {
		t.Errorf("output missing slide title: %s", output.Content)
	}

	// Suppress unused variable warning
	_ = listOutput
}

func TestHandleGetSlideNotFound(t *testing.T) {
	content := `---
marp: true
---

# Only Slide
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	_, _, err := handleGetSlide(ctx, nil, GetSlideInput{
		Path:    path,
		SlideID: "nonexistent",
	})
	if err == nil {
		t.Error("expected error for nonexistent slide")
	}
}

func TestHandlePlanChanges(t *testing.T) {
	content := `---
marp: true
---

# Original Title
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	_, output, err := handlePlanChanges(ctx, nil, PlanChangesInput{
		Path: path,
		Desired: model.Deck{
			Title: "New Title",
		},
	})
	if err != nil {
		t.Fatalf("handlePlanChanges failed: %v", err)
	}

	// Should have changes
	if output.IsEmpty {
		t.Error("expected non-empty diff for title change")
	}

	if !strings.Contains(output.Content, "plan deck") {
		t.Errorf("output missing plan header: %s", output.Content)
	}
}

func TestHandleApplyChangesWithoutConfirm(t *testing.T) {
	content := `---
marp: true
---

# Test
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	diff := model.Diff{
		DeckID: "test",
		Changes: []model.Change{
			{Op: model.ChangeUpdate, Path: "title", OldValue: "Test", NewValue: "New"},
		},
	}

	_, output, err := handleApplyChanges(ctx, nil, ApplyChangesInput{
		Path:    path,
		Diff:    diff,
		Confirm: false,
	})
	if err != nil {
		t.Fatalf("handleApplyChanges failed: %v", err)
	}

	if output.Applied {
		t.Error("should not be applied without confirm")
	}
	if output.Message == "" {
		t.Error("expected message about confirmation")
	}
}

func TestHandleApplyChangesWithConfirm(t *testing.T) {
	content := `---
marp: true
---

# Original
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	diff := model.Diff{
		DeckID: "test",
		Changes: []model.Change{
			{Op: model.ChangeUpdate, Path: "title", OldValue: "Original", NewValue: "New"},
		},
	}

	_, output, err := handleApplyChanges(ctx, nil, ApplyChangesInput{
		Path:    path,
		Diff:    diff,
		Confirm: true,
	})
	if err != nil {
		t.Fatalf("handleApplyChanges failed: %v", err)
	}

	if !output.Applied {
		t.Error("expected changes to be applied")
	}
}

func TestHandleCreateDeck(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "new.md")
	ctx := context.Background()

	deck := model.Deck{
		Title: "Created Deck",
		Sections: []model.Section{
			{
				ID: "main",
				Slides: []model.Slide{
					{ID: "s1", Title: "First Slide"},
				},
			},
		},
	}

	_, output, err := handleCreateDeck(ctx, nil, CreateDeckInput{
		Path: path,
		Deck: deck,
	})
	if err != nil {
		t.Fatalf("handleCreateDeck failed: %v", err)
	}

	if output.Path == "" {
		t.Error("expected non-empty path")
	}

	// Verify file exists
	if _, err := os.Stat(output.Path); os.IsNotExist(err) {
		t.Error("created file does not exist")
	}
}

func TestHandleUpdateSlideWithoutConfirm(t *testing.T) {
	content := `---
marp: true
---

# Slide to Update

Original content.
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	// Get slide ID
	listResult, err := ops.ListSlidesFromPath(ctx, path, "toon")
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	slideID := listResult.Slides[0].ID

	_, output, err := handleUpdateSlide(ctx, nil, UpdateSlideInput{
		Path:    path,
		SlideID: slideID,
		Updates: model.Slide{Title: "Updated Title"},
		Confirm: false,
	})
	if err != nil {
		t.Fatalf("handleUpdateSlide failed: %v", err)
	}

	if output.Updated {
		t.Error("should not be updated without confirm")
	}
}

func TestHandleUpdateSlideWithConfirm(t *testing.T) {
	content := `---
marp: true
---

# Original Title

Some content.
`
	path := createTestPresentation(t, content)
	ctx := context.Background()

	// Get slide ID
	listResult, err := ops.ListSlidesFromPath(ctx, path, "toon")
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	slideID := listResult.Slides[0].ID

	_, output, err := handleUpdateSlide(ctx, nil, UpdateSlideInput{
		Path:    path,
		SlideID: slideID,
		Updates: model.Slide{Title: "New Title"},
		Confirm: true,
	})
	if err != nil {
		t.Fatalf("handleUpdateSlide failed: %v", err)
	}

	if !output.Updated {
		t.Error("expected slide to be updated")
	}
}
