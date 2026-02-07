package ops

import (
	"context"

	"github.com/grokify/slidekit/model"
)

// CreateOptions configures the CreateDeck operation.
type CreateOptions struct {
	Backend string
	Path    string
}

// CreateResult contains the result of a CreateDeck operation.
type CreateResult struct {
	Ref     model.Ref
	Message string
}

// CreateDeck creates a new presentation.
func CreateDeck(ctx context.Context, deck *model.Deck, opts CreateOptions) (*CreateResult, error) {
	backendName := opts.Backend
	if backendName == "" {
		backendName = "marp"
	}

	backend, err := DefaultRegistry.Get(backendName)
	if err != nil {
		return nil, err
	}

	// If a path is specified, use it as the deck ID for file naming
	if opts.Path != "" && deck.ID == "" {
		// Strip .md extension for ID
		id := opts.Path
		if len(id) > 3 && id[len(id)-3:] == ".md" {
			id = id[:len(id)-3]
		}
		deck.ID = id
	}

	ref, err := backend.Create(ctx, deck)
	if err != nil {
		return nil, err
	}

	return &CreateResult{
		Ref:     ref,
		Message: "Deck created at " + ref.Path,
	}, nil
}
