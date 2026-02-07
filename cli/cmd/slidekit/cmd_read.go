package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/ops"
)

var readFormat string

var readCmd = &cobra.Command{
	Use:   "read <file>",
	Short: "Read a presentation file",
	Long: `Read a presentation file and output its content in TOON (default) or JSON format.

TOON format is optimized for AI agents, being ~8x more token-efficient than JSON.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Validate format
		f := format.Format(readFormat)
		if readFormat != "" && !f.IsValid() {
			return fmt.Errorf("invalid format: %s (use 'toon' or 'json')", readFormat)
		}

		result, err := ops.ReadDeckFromPath(context.Background(), path, ops.ReadOptions{
			Format: f,
		})
		if err != nil {
			return fmt.Errorf("reading deck: %w", err)
		}

		fmt.Fprint(os.Stdout, result.Output)
		return nil
	},
}

func init() {
	readCmd.Flags().StringVarP(&readFormat, "format", "f", "toon", "Output format: toon or json")
}
