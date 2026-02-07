package tools

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

// ApplyChangesInput is the input for the apply_changes tool.
type ApplyChangesInput struct {
	Path    string     `json:"path" jsonschema:"description=path to the presentation file"`
	Diff    model.Diff `json:"diff" jsonschema:"description=the diff to apply"`
	Confirm bool       `json:"confirm" jsonschema:"description=must be true to actually apply changes"`
}

// ApplyChangesOutput is the output for the apply_changes tool.
type ApplyChangesOutput struct {
	Applied bool   `json:"applied" jsonschema:"description=true if changes were applied"`
	Message string `json:"message" jsonschema:"description=status message"`
}

var applyChangesTool = &mcp.Tool{
	Name:        "apply_changes",
	Description: "Apply a diff to a presentation. Requires confirm=true to make changes.",
}

func handleApplyChanges(ctx context.Context, req *mcp.CallToolRequest, input ApplyChangesInput) (*mcp.CallToolResult, ApplyChangesOutput, error) {
	result, err := ops.ApplyChangesFromPath(ctx, input.Path, &input.Diff, ops.ApplyOptions{
		Confirm: input.Confirm,
	})
	if err != nil {
		if errors.Is(err, ops.ErrConfirmRequired) {
			return nil, ApplyChangesOutput{
				Applied: false,
				Message: result.Message,
			}, nil
		}
		return nil, ApplyChangesOutput{}, err
	}

	return nil, ApplyChangesOutput{
		Applied: result.Applied,
		Message: result.Message,
	}, nil
}
