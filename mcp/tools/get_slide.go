package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/ops"
)

// GetSlideInput is the input for the get_slide tool.
type GetSlideInput struct {
	Path    string `json:"path" jsonschema:"description=path to the presentation file"`
	SlideID string `json:"slide_id" jsonschema:"description=ID of the slide to retrieve"`
	Format  string `json:"format,omitempty" jsonschema:"description=output format: toon (default) or json"`
}

// GetSlideOutput is the output for the get_slide tool.
type GetSlideOutput struct {
	Content string `json:"content" jsonschema:"description=the slide content in the requested format"`
}

var getSlideTool = &mcp.Tool{
	Name:        "get_slide",
	Description: "Get a single slide by ID from a presentation",
}

func handleGetSlide(ctx context.Context, req *mcp.CallToolRequest, input GetSlideInput) (*mcp.CallToolResult, GetSlideOutput, error) {
	f := format.Format(input.Format)
	if input.Format == "" {
		f = format.FormatTOON
	}

	result, err := ops.GetSlideFromPath(ctx, input.Path, input.SlideID, f)
	if err != nil {
		return nil, GetSlideOutput{}, err
	}

	return nil, GetSlideOutput{Content: result.Output}, nil
}
