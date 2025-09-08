package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initTool1() []server.ServerTool {
	tools := []server.ServerTool{
		{Tool: mcp.NewTool("list_browsers",
			mcp.WithDescription("List the available browsers"),
		), Handler: s.tool1},
	}
	return tools
}

func (s *Server) tool1(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return NewTextResult("Chrome, Firefox, Safari", nil), nil
}
