package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/ops"
)

// ListSlidesInput is the input for the list_slides tool.
type ListSlidesInput struct {
	Path   string `json:"path" jsonschema:"description=path to the presentation file"`
	Format string `json:"format,omitempty" jsonschema:"description=output format: toon (default) or json"`
}

// ListSlidesOutput is the output for the list_slides tool.
type ListSlidesOutput struct {
	Content string `json:"content" jsonschema:"description=list of slides with IDs and titles"`
}

var listSlidesTool = &mcp.Tool{
	Name:        "list_slides",
	Description: "List all slides in a presentation with their IDs and titles",
}

func handleListSlides(ctx context.Context, req *mcp.CallToolRequest, input ListSlidesInput) (*mcp.CallToolResult, ListSlidesOutput, error) {
	f := format.Format(input.Format)
	if input.Format == "" {
		f = format.FormatTOON
	}

	result, err := ops.ListSlidesFromPath(ctx, input.Path, f)
	if err != nil {
		return nil, ListSlidesOutput{}, err
	}

	return nil, ListSlidesOutput{Content: result.Output}, nil
}
