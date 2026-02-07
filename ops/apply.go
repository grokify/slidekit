package ops

import (
	"context"
	"errors"

	"github.com/grokify/slidekit/model"
)

// ErrConfirmRequired is returned when confirm=false.
var ErrConfirmRequired = errors.New("confirmation required: set confirm=true to apply changes")

// ApplyOptions configures the ApplyChanges operation.
type ApplyOptions struct {
	Confirm bool
}

// ApplyResult contains the result of an ApplyChanges operation.
type ApplyResult struct {
	Applied bool
	Message string
}

// ApplyChanges applies a diff to the presentation.
func ApplyChanges(ctx context.Context, ref model.Ref, diff *model.Diff, opts ApplyOptions) (*ApplyResult, error) {
	if !opts.Confirm {
		return &ApplyResult{
			Applied: false,
			Message: "Set confirm=true to apply changes",
		}, ErrConfirmRequired
	}

	if diff.IsEmpty() {
		return &ApplyResult{
			Applied: false,
			Message: "No changes to apply",
		}, nil
	}

	backend, err := DefaultRegistry.Get(ref.Backend)
	if err != nil {
		return nil, err
	}

	if err := backend.Apply(ctx, ref, diff); err != nil {
		return nil, err
	}

	return &ApplyResult{
		Applied: true,
		Message: "Changes applied successfully",
	}, nil
}

// ApplyChangesFromPath is a convenience function that detects the backend.
func ApplyChangesFromPath(ctx context.Context, path string, diff *model.Diff, opts ApplyOptions) (*ApplyResult, error) {
	backendName := DetectBackend(path)
	ref := model.Ref{
		Backend: backendName,
		Path:    path,
	}
	return ApplyChanges(ctx, ref, diff, opts)
}
