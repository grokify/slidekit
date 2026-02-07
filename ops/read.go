package ops

import (
	"context"
	"encoding/json"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
)

// ReadOptions configures the ReadDeck operation.
type ReadOptions struct {
	Format format.Format
}

// ReadResult contains the result of a ReadDeck operation.
type ReadResult struct {
	Deck   *model.Deck
	Output string
}

// ReadDeck reads a presentation and returns it in the requested format.
func ReadDeck(ctx context.Context, ref model.Ref, opts ReadOptions) (*ReadResult, error) {
	backend, err := DefaultRegistry.Get(ref.Backend)
	if err != nil {
		return nil, err
	}

	deck, err := backend.Read(ctx, ref)
	if err != nil {
		return nil, err
	}

	output, err := encodeDeck(deck, opts.Format)
	if err != nil {
		return nil, err
	}

	return &ReadResult{
		Deck:   deck,
		Output: output,
	}, nil
}

// ReadDeckFromPath is a convenience function that detects the backend.
func ReadDeckFromPath(ctx context.Context, path string, opts ReadOptions) (*ReadResult, error) {
	backendName := DetectBackend(path)
	ref := model.Ref{
		Backend: backendName,
		Path:    path,
	}
	return ReadDeck(ctx, ref, opts)
}

// encodeDeck serializes a deck to the requested format.
func encodeDeck(deck *model.Deck, f format.Format) (string, error) {
	switch f {
	case format.FormatJSON:
		data, err := json.MarshalIndent(deck, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case format.FormatTOON, "":
		encoder := format.NewTOONEncoder()
		return encoder.EncodeDeck(deck), nil
	default:
		encoder := format.NewTOONEncoder()
		return encoder.EncodeDeck(deck), nil
	}
}
