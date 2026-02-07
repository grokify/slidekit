package ops

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/slidekit/model"
)

func TestCreateDeck(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	deck := &model.Deck{
		Title: "New Presentation",
		Sections: []model.Section{
			{
				ID:    "intro",
				Title: "Introduction",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitle,
						Title:  "Welcome",
					},
					{
						ID:     "s2",
						Layout: model.LayoutTitleBody,
						Title:  "Overview",
						Body: []model.Block{
							model.NewBullet("First point", 0),
							model.NewBullet("Second point", 0),
						},
					},
				},
			},
		},
	}

	result, err := CreateDeck(ctx, deck, CreateOptions{
		Backend: "marp",
		Path:    filepath.Join(tmpDir, "created.md"),
	})
	if err != nil {
		t.Fatalf("CreateDeck failed: %v", err)
	}

	if result.Ref.Path == "" {
		t.Error("expected non-empty path in result")
	}

	// Verify file was created
	content, err := os.ReadFile(result.Ref.Path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	// Check content
	if !strings.Contains(string(content), "marp: true") {
		t.Error("created file missing marp frontmatter")
	}
	if !strings.Contains(string(content), "Welcome") {
		t.Error("created file missing slide title")
	}
	if !strings.Contains(string(content), "First point") {
		t.Error("created file missing bullet content")
	}
}

func TestCreateDeckAutoDetectBackend(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	deck := &model.Deck{
		Title: "Auto Detect Test",
		Sections: []model.Section{
			{
				ID: "main",
				Slides: []model.Slide{
					{ID: "s1", Title: "Slide One"},
				},
			},
		},
	}

	// Don't specify backend, let it auto-detect from .md extension
	result, err := CreateDeck(ctx, deck, CreateOptions{
		Path: filepath.Join(tmpDir, "auto.md"),
	})
	if err != nil {
		t.Fatalf("CreateDeck failed: %v", err)
	}

	if result.Ref.Backend != "marp" {
		t.Errorf("expected backend 'marp', got %q", result.Ref.Backend)
	}
}

func TestCreateDeckWithID(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	deck := &model.Deck{
		ID:    "my-presentation",
		Title: "ID Test",
	}

	// Create without specifying path - should use ID
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Logf("failed to restore dir: %v", err)
		}
	}()

	result, err := CreateDeck(ctx, deck, CreateOptions{Backend: "marp"})
	if err != nil {
		t.Fatalf("CreateDeck failed: %v", err)
	}

	if !strings.Contains(result.Ref.Path, "my-presentation") {
		t.Errorf("expected path to contain 'my-presentation', got %q", result.Ref.Path)
	}
}

func TestCreateDeckUnknownBackend(t *testing.T) {
	ctx := context.Background()
	deck := &model.Deck{Title: "Test"}

	_, err := CreateDeck(ctx, deck, CreateOptions{Backend: "unknown"})
	if err == nil {
		t.Error("expected error for unknown backend")
	}
}

func TestCreateDeckMessage(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	deck := &model.Deck{Title: "Message Test"}

	result, err := CreateDeck(ctx, deck, CreateOptions{
		Backend: "marp",
		Path:    filepath.Join(tmpDir, "msg.md"),
	})
	if err != nil {
		t.Fatalf("CreateDeck failed: %v", err)
	}

	if !strings.Contains(result.Message, "created") {
		t.Errorf("expected message to mention 'created', got: %s", result.Message)
	}
}
