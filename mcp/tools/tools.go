// Package tools provides MCP tool implementations for slidekit.
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all slidekit tools with the MCP server.
func RegisterAll(srv *mcp.Server) {
	mcp.AddTool(srv, readDeckTool, handleReadDeck)
	mcp.AddTool(srv, listSlidesTool, handleListSlides)
	mcp.AddTool(srv, getSlideTool, handleGetSlide)
	mcp.AddTool(srv, planChangesTool, handlePlanChanges)
	mcp.AddTool(srv, applyChangesTool, handleApplyChanges)
	mcp.AddTool(srv, createDeckTool, handleCreateDeck)
	mcp.AddTool(srv, updateSlideTool, handleUpdateSlide)
}
