package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

var (
	applyDiff    string
	applyConfirm bool
)

var applyCmd = &cobra.Command{
	Use:   "apply <file>",
	Short: "Apply changes to a presentation",
	Long: `Apply changes from a diff file to a presentation.

The diff must be provided as a JSON file using the --diff flag.
The --confirm flag is required to actually apply changes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if applyDiff == "" {
			return fmt.Errorf("--diff flag is required")
		}

		// Read diff from file
		diffData, err := os.ReadFile(applyDiff)
		if err != nil {
			return fmt.Errorf("reading diff file: %w", err)
		}

		var diff model.Diff
		if err := json.Unmarshal(diffData, &diff); err != nil {
			return fmt.Errorf("parsing diff file: %w", err)
		}

		result, err := ops.ApplyChangesFromPath(context.Background(), path, &diff, ops.ApplyOptions{
			Confirm: applyConfirm,
		})
		if err != nil {
			if errors.Is(err, ops.ErrConfirmRequired) {
				fmt.Println(result.Message)
				fmt.Println("Use --confirm to apply changes")
				return nil
			}
			return fmt.Errorf("applying changes: %w", err)
		}

		fmt.Println(result.Message)
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyDiff, "diff", "d", "", "Path to diff file (JSON)")
	applyCmd.Flags().BoolVar(&applyConfirm, "confirm", false, "Confirm application of changes")
}
