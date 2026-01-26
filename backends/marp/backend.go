package marp

import (
	"context"
	"fmt"

	"github.com/grokify/slidekit/model"
)

// Backend implements the model.Backend interface for Marp Markdown files.
type Backend struct {
	reader *Reader
	writer *Writer
}

// NewBackend creates a new Marp backend.
func NewBackend() *Backend {
	return &Backend{
		reader: NewReader(),
		writer: NewWriter(),
	}
}

// Info returns backend metadata.
func (b *Backend) Info() model.BackendInfo {
	return model.BackendInfo{
		Name:    "marp",
		Version: "0.1.0",
		Capabilities: []string{
			model.CapabilityRead,
			model.CapabilityWrite,
			model.CapabilityCreate,
			model.CapabilitySections,
		},
	}
}

// Read loads a Marp presentation from a file.
func (b *Backend) Read(ctx context.Context, ref model.Ref) (*model.Deck, error) {
	return b.reader.Read(ctx, ref)
}

// Plan computes changes needed to reach desired state.
func (b *Backend) Plan(_ context.Context, ref model.Ref, desired *model.Deck) (*model.Diff, error) {
	// Read current state
	current, err := b.reader.ReadFile(ref.Path)
	if err != nil {
		return nil, fmt.Errorf("reading current deck: %w", err)
	}

	return computeDiff(current, desired), nil
}

// Apply writes the diff to the file. For Marp, this means rewriting the file.
func (b *Backend) Apply(_ context.Context, ref model.Ref, diff *model.Diff) error {
	if diff.IsEmpty() {
		return nil
	}

	// For Marp files, we regenerate the whole file from the desired state.
	// A more sophisticated implementation could do surgical edits.
	current, err := b.reader.ReadFile(ref.Path)
	if err != nil {
		return fmt.Errorf("reading current deck: %w", err)
	}

	applyDiff(current, diff)
	return b.writer.WriteFile(current, ref.Path)
}

// Create creates a new Marp presentation file.
func (b *Backend) Create(_ context.Context, deck *model.Deck) (model.Ref, error) {
	// Default path if not set
	path := "presentation.md"
	if deck.ID != "" {
		path = deck.ID + ".md"
	}

	if err := b.writer.WriteFile(deck, path); err != nil {
		return model.Ref{}, fmt.Errorf("writing deck: %w", err)
	}

	return model.Ref{
		Backend: "marp",
		Path:    path,
	}, nil
}

// computeDiff compares two decks and produces a diff.
func computeDiff(current, desired *model.Deck) *model.Diff {
	diff := model.NewDiff(current.ID)

	// Compare title
	if current.Title != desired.Title {
		diff.AddChange(model.NewUpdateChange("title", current.Title, desired.Title))
	}

	// Compare sections
	currentSections := make(map[string]*model.Section)
	for i := range current.Sections {
		currentSections[current.Sections[i].ID] = &current.Sections[i]
	}

	for _, ds := range desired.Sections {
		cs, exists := currentSections[ds.ID]
		if !exists {
			diff.AddChange(model.NewAddChange("sections/"+ds.ID, ds))
			continue
		}
		// Compare slides within section
		compareSectionSlides(diff, cs, &ds)
	}

	// Check for removed sections
	desiredSections := make(map[string]bool)
	for _, s := range desired.Sections {
		desiredSections[s.ID] = true
	}
	for _, s := range current.Sections {
		if !desiredSections[s.ID] {
			diff.AddChange(model.NewRemoveChange("sections/"+s.ID, s))
		}
	}

	return diff
}

// compareSectionSlides compares slides between two sections.
func compareSectionSlides(diff *model.Diff, current, desired *model.Section) {
	currentSlides := make(map[string]*model.Slide)
	for i := range current.Slides {
		currentSlides[current.Slides[i].ID] = &current.Slides[i]
	}

	for _, ds := range desired.Slides {
		cs, exists := currentSlides[ds.ID]
		if !exists {
			diff.AddChange(model.NewAddChange(
				fmt.Sprintf("sections/%s/slides/%s", current.ID, ds.ID), ds))
			continue
		}
		// Compare slide content
		if cs.Title != ds.Title {
			diff.AddChange(model.NewUpdateChange(
				fmt.Sprintf("sections/%s/slides/%s/title", current.ID, ds.ID),
				cs.Title, ds.Title))
		}
	}
}

// applyDiff applies changes from a diff to a deck (in-place).
func applyDiff(deck *model.Deck, diff *model.Diff) {
	for _, change := range diff.Changes {
		switch change.Op {
		case model.ChangeUpdate:
			if change.Path == "title" {
				if v, ok := change.NewValue.(string); ok {
					deck.Title = v
				}
			}
			// Additional path-based updates would go here
		case model.ChangeAdd:
			// Handle additions
		case model.ChangeRemove:
			// Handle removals
		}
	}
}
