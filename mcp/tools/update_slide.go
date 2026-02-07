package tools

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

// UpdateSlideInput is the input for the update_slide tool.
type UpdateSlideInput struct {
	Path    string      `json:"path" jsonschema:"description=path to the presentation file"`
	SlideID string      `json:"slide_id" jsonschema:"description=ID of the slide to update"`
	Updates model.Slide `json:"updates" jsonschema:"description=fields to update (title, subtitle, body, notes)"`
	Confirm bool        `json:"confirm" jsonschema:"description=must be true to actually apply changes"`
}

// UpdateSlideOutput is the output for the update_slide tool.
type UpdateSlideOutput struct {
	Updated bool   `json:"updated" jsonschema:"description=true if the slide was updated"`
	Message string `json:"message" jsonschema:"description=status message"`
}

var updateSlideTool = &mcp.Tool{
	Name:        "update_slide",
	Description: "Update a single slide. Requires confirm=true to make changes.",
}

func handleUpdateSlide(ctx context.Context, req *mcp.CallToolRequest, input UpdateSlideInput) (*mcp.CallToolResult, UpdateSlideOutput, error) {
	backendName := ops.DetectBackend(input.Path)
	ref := model.Ref{
		Backend: backendName,
		Path:    input.Path,
	}

	result, err := ops.UpdateSlide(ctx, ref, input.SlideID, &input.Updates, ops.UpdateSlideOptions{
		Confirm: input.Confirm,
	})
	if err != nil {
		if errors.Is(err, ops.ErrConfirmRequired) {
			return nil, UpdateSlideOutput{
				Updated: false,
				Message: result.Message,
			}, nil
		}
		return nil, UpdateSlideOutput{}, err
	}

	return nil, UpdateSlideOutput{
		Updated: result.Updated,
		Message: result.Message,
	}, nil
}
