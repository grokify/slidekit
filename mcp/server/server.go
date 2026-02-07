// Package server provides the MCP server for slidekit.
package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/slidekit/mcp/tools"
)

// Server wraps the MCP server with slidekit-specific configuration.
type Server struct {
	server *mcp.Server
}

// New creates a new slidekit MCP server.
func New(version string) *Server {
	impl := &mcp.Implementation{
		Name:    "slidekit",
		Version: version,
	}

	srv := mcp.NewServer(impl, nil)

	// Register all tools
	tools.RegisterAll(srv)

	return &Server{server: srv}
}

// Run starts the server on stdio transport.
func (s *Server) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}
