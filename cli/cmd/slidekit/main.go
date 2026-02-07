// Package main provides the slidekit CLI for managing presentations.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/backends/marp"
	"github.com/grokify/slidekit/ops"
)

var (
	// Version is set at build time.
	Version = "dev"
)

func main() {
	// Register backends
	ops.DefaultRegistry.Register("marp", marp.NewBackend())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "slidekit",
	Short: "A toolkit for managing presentations",
	Long: `slidekit is a CLI for reading, planning, and modifying presentations.

Supports multiple backends including Marp Markdown.`,
	Version: Version,
}

func init() {
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(serveCmd)
}
