package ops

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
)

func TestListSlides(t *testing.T) {
	content := `---
marp: true
---

# First Slide

---

## Second Slide

Content here.

---

## Third Slide

More content.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "list.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	result, err := ListSlides(ctx, ref, format.FormatTOON)
	if err != nil {
		t.Fatalf("ListSlides failed: %v", err)
	}

	if len(result.Slides) != 3 {
		t.Errorf("expected 3 slides, got %d", len(result.Slides))
	}

	// Check first slide
	if result.Slides[0].Title != "First Slide" {
		t.Errorf("expected first slide title 'First Slide', got %q", result.Slides[0].Title)
	}

	// Check output format
	if !strings.Contains(result.Output, "slide") {
		t.Errorf("output missing slide info: %s", result.Output)
	}
}

func TestListSlidesFromPath(t *testing.T) {
	content := `---
marp: true
---

# Only Slide
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "single.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	result, err := ListSlidesFromPath(ctx, path, format.FormatJSON)
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	if len(result.Slides) != 1 {
		t.Errorf("expected 1 slide, got %d", len(result.Slides))
	}

	// JSON format check
	if !strings.Contains(result.Output, `"title"`) {
		t.Errorf("JSON output missing title field: %s", result.Output)
	}
}

func TestGetSlide(t *testing.T) {
	content := `---
marp: true
---

# Title Slide

---

## Target Slide

- Bullet one
- Bullet two

---

## Last Slide
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "get.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()

	// First, list slides to get IDs
	listResult, err := ListSlidesFromPath(ctx, path, format.FormatTOON)
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	if len(listResult.Slides) < 2 {
		t.Fatalf("expected at least 2 slides, got %d", len(listResult.Slides))
	}

	// Get the second slide by ID
	slideID := listResult.Slides[1].ID
	ref := model.Ref{Backend: "marp", Path: path}

	result, err := GetSlide(ctx, ref, slideID, format.FormatTOON)
	if err != nil {
		t.Fatalf("GetSlide failed: %v", err)
	}

	if result.Slide == nil {
		t.Fatal("expected non-nil slide")
	}
	if result.Slide.Title != "Target Slide" {
		t.Errorf("expected title 'Target Slide', got %q", result.Slide.Title)
	}
	if len(result.Slide.Body) != 2 {
		t.Errorf("expected 2 body blocks, got %d", len(result.Slide.Body))
	}
}

func TestGetSlideFromPath(t *testing.T) {
	content := `---
marp: true
---

# The Slide

Content here.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "getpath.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()

	// List to get ID
	listResult, err := ListSlidesFromPath(ctx, path, format.FormatTOON)
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	slideID := listResult.Slides[0].ID
	result, err := GetSlideFromPath(ctx, path, slideID, format.FormatJSON)
	if err != nil {
		t.Fatalf("GetSlideFromPath failed: %v", err)
	}

	if !strings.Contains(result.Output, `"title": "The Slide"`) {
		t.Errorf("JSON output missing expected title: %s", result.Output)
	}
}

func TestGetSlideNotFound(t *testing.T) {
	content := `---
marp: true
---

# Only Slide
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "notfound.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	_, err := GetSlide(ctx, ref, "nonexistent-id", format.FormatTOON)
	if err == nil {
		t.Error("expected error for nonexistent slide")
	}
	if !strings.Contains(err.Error(), "slide not found") {
		t.Errorf("expected 'slide not found' error, got: %v", err)
	}
}

func TestSlideInfoFields(t *testing.T) {
	content := `---
marp: true
---

<!-- _class: section-divider -->

# Section Title

---

## Regular Slide

Content.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "info.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	result, err := ListSlidesFromPath(ctx, path, format.FormatTOON)
	if err != nil {
		t.Fatalf("ListSlidesFromPath failed: %v", err)
	}

	// Check that SlideInfo has correct fields
	for _, slide := range result.Slides {
		if slide.ID == "" {
			t.Error("slide ID should not be empty")
		}
		if slide.SectionID == "" {
			t.Error("slide section ID should not be empty")
		}
		if slide.Layout == "" {
			t.Error("slide layout should not be empty")
		}
	}

	// First slide should be section layout
	if result.Slides[0].Layout != string(model.LayoutSection) {
		t.Errorf("expected section layout, got %q", result.Slides[0].Layout)
	}
}
