package ops

import (
	"context"
	"encoding/json"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
)

// PlanOptions configures the PlanChanges operation.
type PlanOptions struct {
	Format format.Format
}

// PlanResult contains the result of a PlanChanges operation.
type PlanResult struct {
	Diff   *model.Diff
	Output string
}

// PlanChanges computes the diff between current and desired states.
func PlanChanges(ctx context.Context, ref model.Ref, desired *model.Deck, opts PlanOptions) (*PlanResult, error) {
	backend, err := DefaultRegistry.Get(ref.Backend)
	if err != nil {
		return nil, err
	}

	diff, err := backend.Plan(ctx, ref, desired)
	if err != nil {
		return nil, err
	}

	output, err := encodeDiff(diff, opts.Format)
	if err != nil {
		return nil, err
	}

	return &PlanResult{
		Diff:   diff,
		Output: output,
	}, nil
}

// PlanChangesFromPath is a convenience function that detects the backend.
func PlanChangesFromPath(ctx context.Context, path string, desired *model.Deck, opts PlanOptions) (*PlanResult, error) {
	backendName := DetectBackend(path)
	ref := model.Ref{
		Backend: backendName,
		Path:    path,
	}
	return PlanChanges(ctx, ref, desired, opts)
}

// encodeDiff serializes a diff to the requested format.
func encodeDiff(diff *model.Diff, f format.Format) (string, error) {
	switch f {
	case format.FormatJSON:
		data, err := json.MarshalIndent(diff, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case format.FormatTOON, "":
		encoder := format.NewTOONEncoder()
		return encoder.EncodeDiff(diff), nil
	default:
		encoder := format.NewTOONEncoder()
		return encoder.EncodeDiff(diff), nil
	}
}
