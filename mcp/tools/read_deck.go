package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/ops"
)

// ReadDeckInput is the input for the read_deck tool.
type ReadDeckInput struct {
	Path   string `json:"path" jsonschema:"description=path to the presentation file"`
	Format string `json:"format,omitempty" jsonschema:"description=output format: toon (default) or json"`
}

// ReadDeckOutput is the output for the read_deck tool.
type ReadDeckOutput struct {
	Content string `json:"content" jsonschema:"description=the presentation content in the requested format"`
}

var readDeckTool = &mcp.Tool{
	Name:        "read_deck",
	Description: "Read a presentation file and return its content in TOON (default) or JSON format. TOON is ~8x more token-efficient than JSON.",
}

func handleReadDeck(ctx context.Context, req *mcp.CallToolRequest, input ReadDeckInput) (*mcp.CallToolResult, ReadDeckOutput, error) {
	f := format.Format(input.Format)
	if input.Format == "" {
		f = format.FormatTOON
	}

	result, err := ops.ReadDeckFromPath(ctx, input.Path, ops.ReadOptions{
		Format: f,
	})
	if err != nil {
		return nil, ReadDeckOutput{}, err
	}

	return nil, ReadDeckOutput{Content: result.Output}, nil
}
