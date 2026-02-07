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

func TestPlanChanges(t *testing.T) {
	content := `---
marp: true
---

# Original Title

---

## Slide Two

Original content.
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "plan.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	// Create desired state with changed title
	desired := &model.Deck{
		Title: "New Title",
		Sections: []model.Section{
			{
				ID: "section-0",
				Slides: []model.Slide{
					{ID: "s0-0", Title: "New Title"},
					{ID: "s0-1", Title: "Slide Two"},
				},
			},
		},
	}

	result, err := PlanChanges(ctx, ref, desired, PlanOptions{Format: format.FormatTOON})
	if err != nil {
		t.Fatalf("PlanChanges failed: %v", err)
	}

	if result.Diff == nil {
		t.Fatal("expected non-nil diff")
	}

	// Should have at least one change (title)
	if result.Diff.IsEmpty() {
		t.Error("expected non-empty diff for title change")
	}

	// Check TOON output
	if !strings.Contains(result.Output, "plan deck") {
		t.Errorf("TOON output missing plan header: %s", result.Output)
	}
}

func TestPlanChangesNoChanges(t *testing.T) {
	content := `---
marp: true
---

# Same Title
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nochange.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()

	// Read current state first
	readResult, err := ReadDeckFromPath(ctx, path, ReadOptions{})
	if err != nil {
		t.Fatalf("ReadDeckFromPath failed: %v", err)
	}

	// Use same state as desired
	result, err := PlanChangesFromPath(ctx, path, readResult.Deck, PlanOptions{})
	if err != nil {
		t.Fatalf("PlanChangesFromPath failed: %v", err)
	}

	if !result.Diff.IsEmpty() {
		t.Errorf("expected empty diff when no changes, got %d changes", result.Diff.ChangeCount())
	}
}

func TestPlanChangesJSON(t *testing.T) {
	content := `---
marp: true
---

# Test
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "planjson.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	ref := model.Ref{Backend: "marp", Path: path}

	desired := &model.Deck{Title: "Changed"}

	result, err := PlanChanges(ctx, ref, desired, PlanOptions{Format: format.FormatJSON})
	if err != nil {
		t.Fatalf("PlanChanges failed: %v", err)
	}

	if !strings.Contains(result.Output, `"changes"`) {
		t.Errorf("JSON output missing changes field: %s", result.Output)
	}
}

func TestPlanChangesFromPath(t *testing.T) {
	content := `---
marp: true
---

# Original
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "planpath.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ctx := context.Background()
	desired := &model.Deck{Title: "Modified"}

	result, err := PlanChangesFromPath(ctx, path, desired, PlanOptions{})
	if err != nil {
		t.Fatalf("PlanChangesFromPath failed: %v", err)
	}

	if result.Diff.IsEmpty() {
		t.Error("expected non-empty diff")
	}
}
