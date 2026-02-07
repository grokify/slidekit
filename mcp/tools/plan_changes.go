package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

// PlanChangesInput is the input for the plan_changes tool.
type PlanChangesInput struct {
	Path    string     `json:"path" jsonschema:"description=path to the presentation file"`
	Desired model.Deck `json:"desired" jsonschema:"description=the desired state of the deck"`
	Format  string     `json:"format,omitempty" jsonschema:"description=output format: toon (default) or json"`
}

// PlanChangesOutput is the output for the plan_changes tool.
type PlanChangesOutput struct {
	Content string `json:"content" jsonschema:"description=the diff between current and desired state"`
	IsEmpty bool   `json:"is_empty" jsonschema:"description=true if there are no changes"`
}

var planChangesTool = &mcp.Tool{
	Name:        "plan_changes",
	Description: "Compute the diff between the current presentation and a desired state",
}

func handlePlanChanges(ctx context.Context, req *mcp.CallToolRequest, input PlanChangesInput) (*mcp.CallToolResult, PlanChangesOutput, error) {
	f := format.Format(input.Format)
	if input.Format == "" {
		f = format.FormatTOON
	}

	result, err := ops.PlanChangesFromPath(ctx, input.Path, &input.Desired, ops.PlanOptions{
		Format: f,
	})
	if err != nil {
		return nil, PlanChangesOutput{}, err
	}

	return nil, PlanChangesOutput{
		Content: result.Output,
		IsEmpty: result.Diff.IsEmpty(),
	}, nil
}

// PlanChangesInputJSON is an alternative input that accepts JSON for desired state.
type PlanChangesInputJSON struct {
	Path        string          `json:"path" jsonschema:"description=path to the presentation file"`
	DesiredJSON json.RawMessage `json:"desired_json" jsonschema:"description=the desired state as JSON"`
	Format      string          `json:"format,omitempty" jsonschema:"description=output format: toon (default) or json"`
}
