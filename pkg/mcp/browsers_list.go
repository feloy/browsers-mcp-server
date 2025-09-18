package mcp

import (
	"context"
	"strings"

	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListBrowsersResult struct {
	Names []string `json:"browsers"`
}

func (s *Server) initBrowsersList() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_browsers",
				mcp.WithDescription("List the available browsers"),
			),
			Handler: s.listBrowsers,
		},
	}
	return tools
}

func (s *Server) listBrowsers(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	browsers := browsers.GetBrowsers()
	names := []string{}
	for _, browser := range browsers {
		names = append(names, browser.Name())
	}
	return NewStructuredResult(strings.Join(names, ", "), ListBrowsersResult{Names: names}, nil), nil
}
