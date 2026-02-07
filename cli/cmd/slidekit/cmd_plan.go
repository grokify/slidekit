package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

var (
	planDesired string
	planFormat  string
)

var planCmd = &cobra.Command{
	Use:   "plan <file>",
	Short: "Show changes between current and desired state",
	Long: `Plan computes the diff between the current presentation and a desired state.

The desired state can be provided as a JSON file using the --desired flag.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if planDesired == "" {
			return fmt.Errorf("--desired flag is required")
		}

		// Read desired state from file
		desiredData, err := os.ReadFile(planDesired)
		if err != nil {
			return fmt.Errorf("reading desired file: %w", err)
		}

		var desired model.Deck
		if err := json.Unmarshal(desiredData, &desired); err != nil {
			return fmt.Errorf("parsing desired file: %w", err)
		}

		// Validate format
		f := format.Format(planFormat)
		if planFormat != "" && !f.IsValid() {
			return fmt.Errorf("invalid format: %s (use 'toon' or 'json')", planFormat)
		}

		result, err := ops.PlanChangesFromPath(context.Background(), path, &desired, ops.PlanOptions{
			Format: f,
		})
		if err != nil {
			return fmt.Errorf("computing plan: %w", err)
		}

		fmt.Fprint(os.Stdout, result.Output)
		return nil
	},
}

func init() {
	planCmd.Flags().StringVarP(&planDesired, "desired", "d", "", "Path to desired state file (JSON)")
	planCmd.Flags().StringVarP(&planFormat, "format", "f", "toon", "Output format: toon or json")
}
