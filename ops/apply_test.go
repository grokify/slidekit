package ops

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/slidekit/model"
)

func TestApplyChangesRequiresConfirm(t *testing.T) {
	content := `---
marp: true
---

# Test
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "apply.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	diff := model.NewDiff("test")
	diff.AddChange(model.NewUpdateChange("title", "Test", "New Title"))

	// Without confirm
	result, err := ApplyChanges(ctx, ref, diff, ApplyOptions{Confirm: false})
	if !errors.Is(err, ErrConfirmRequired) {
		t.Errorf("expected ErrConfirmRequired, got: %v", err)
	}
	if result.Applied {
		t.Error("should not be applied without confirm")
	}
	if result.Message == "" {
		t.Error("expected message about confirmation")
	}
}

func TestApplyChangesWithConfirm(t *testing.T) {
	content := `---
marp: true
---

# Original Title

Some content.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "applyconfirm.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	diff := model.NewDiff("test")
	diff.AddChange(model.NewUpdateChange("title", "Original Title", "New Title"))

	// With confirm
	result, err := ApplyChanges(ctx, ref, diff, ApplyOptions{Confirm: true})
	if err != nil {
		t.Fatalf("ApplyChanges failed: %v", err)
	}
	if !result.Applied {
		t.Error("expected changes to be applied")
	}
}

func TestApplyChangesEmptyDiff(t *testing.T) {
	content := `---
marp: true
---

# Test
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	diff := model.NewDiff("test") // Empty diff

	result, err := ApplyChanges(ctx, ref, diff, ApplyOptions{Confirm: true})
	if err != nil {
		t.Fatalf("ApplyChanges failed: %v", err)
	}
	if result.Applied {
		t.Error("should not report applied for empty diff")
	}
	if result.Message != "No changes to apply" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestApplyChangesFromPath(t *testing.T) {
	content := `---
marp: true
---

# Test
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "applypath.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	diff := model.NewDiff("test")

	// Without confirm
	result, err := ApplyChangesFromPath(ctx, path, diff, ApplyOptions{Confirm: false})
	if !errors.Is(err, ErrConfirmRequired) {
		t.Errorf("expected ErrConfirmRequired, got: %v", err)
	}
	if result.Applied {
		t.Error("should not be applied")
	}
}

func TestApplyChangesUnknownBackend(t *testing.T) {
	ctx := context.Background()
	ref := model.Ref{Backend: "unknown", Path: "/some/path.md"}
	diff := model.NewDiff("test")
	// Add a change so the diff is not empty (empty diffs return early)
	diff.AddChange(model.NewUpdateChange("title", "old", "new"))

	_, err := ApplyChanges(ctx, ref, diff, ApplyOptions{Confirm: true})
	if err == nil {
		t.Error("expected error for unknown backend")
	}
}
