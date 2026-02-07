package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

// CreateDeckInput is the input for the create_deck tool.
type CreateDeckInput struct {
	Path    string     `json:"path" jsonschema:"description=path for the new presentation file"`
	Deck    model.Deck `json:"deck" jsonschema:"description=the deck definition"`
	Backend string     `json:"backend,omitempty" jsonschema:"description=backend to use (default: auto-detect from path)"`
}

// CreateDeckOutput is the output for the create_deck tool.
type CreateDeckOutput struct {
	Path    string `json:"path" jsonschema:"description=path to the created file"`
	Message string `json:"message" jsonschema:"description=status message"`
}

var createDeckTool = &mcp.Tool{
	Name:        "create_deck",
	Description: "Create a new presentation file",
}

func handleCreateDeck(ctx context.Context, req *mcp.CallToolRequest, input CreateDeckInput) (*mcp.CallToolResult, CreateDeckOutput, error) {
	backend := input.Backend
	if backend == "" {
		backend = ops.DetectBackend(input.Path)
	}

	result, err := ops.CreateDeck(ctx, &input.Deck, ops.CreateOptions{
		Backend: backend,
		Path:    input.Path,
	})
	if err != nil {
		return nil, CreateDeckOutput{}, err
	}

	return nil, CreateDeckOutput{
		Path:    result.Ref.Path,
		Message: result.Message,
	}, nil
}
