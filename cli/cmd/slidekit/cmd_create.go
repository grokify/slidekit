package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/model"
	"github.com/grokify/slidekit/ops"
)

var createBackend string

var createCmd = &cobra.Command{
	Use:   "create <file>",
	Short: "Create a new presentation",
	Long: `Create a new presentation from JSON input on stdin.

Example:
  echo '{"title": "My Presentation"}' | slidekit create presentation.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Read deck definition from stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}

		var deck model.Deck
		if err := json.Unmarshal(input, &deck); err != nil {
			return fmt.Errorf("parsing input: %w", err)
		}

		backend := createBackend
		if backend == "" {
			backend = ops.DetectBackend(path)
		}

		result, err := ops.CreateDeck(context.Background(), &deck, ops.CreateOptions{
			Backend: backend,
			Path:    path,
		})
		if err != nil {
			return fmt.Errorf("creating deck: %w", err)
		}

		fmt.Println(result.Message)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createBackend, "backend", "b", "", "Backend to use (default: auto-detect)")
}
