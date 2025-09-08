package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/feloy/mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initTool1() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_browsers",
				mcp.WithDescription("List the available browsers"),
			),
			Handler: s.listBrowsers,
		},
		{
			Tool: mcp.NewTool("list_profiles",
				mcp.WithDescription("List the available profiles for a browser"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the profiles for"),
					mcp.Required(),
				),
			),
			Handler: s.listProfiles,
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
	return NewTextResult(strings.Join(names, ", "), nil), nil
}

func (s *Server) listProfiles(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var browserName string
	ok := false
	if browserName, ok = ctr.GetArguments()["browser"].(string); !ok {
		return NewTextResult("", fmt.Errorf("failed to get list of profiles, missing argument browser")), nil
	}

	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}
	profiles, err := browser.Profiles()
	if err != nil {
		return NewTextResult("", err), nil
	}
	return NewTextResult(strings.Join(profiles, ", "), nil), nil
}
