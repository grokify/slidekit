package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/grokify/slidekit/mcp/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long: `Start the MCP server on stdio transport.

This allows AI assistants to interact with presentations via the
Model Context Protocol (MCP).

Example Claude Code configuration:
  {
    "mcpServers": {
      "slidekit": {
        "command": "/path/to/slidekit",
        "args": ["serve"]
      }
    }
  }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create context that cancels on interrupt
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			cancel()
		}()

		// Create and run server
		srv := server.New(Version)
		return srv.Run(ctx)
	},
}
