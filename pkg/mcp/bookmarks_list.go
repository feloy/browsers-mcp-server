package mcp

import (
	"context"
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initBookmarksList() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_bookmarks",
				mcp.WithDescription("List the available bookmarks for a browser's profile"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the bookmarks for"),
				),
				mcp.WithString(
					"profile",
					mcp.Description("The browser's profile to list the bookmarks for, required if there are multiple profiles"),
				),
			),
			Handler: s.listBookmarks,
		},
	}
	return tools
}

func (s *Server) listBookmarks(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	browserName, err := s.getBrowserName(ctr, "list of bookmarks")
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	profileName, err := s.getProfileName(browser, ctr, "list of bookmarks")
	if err != nil {
		return NewTextResult("", err), nil
	}

	bookmarks, err := browser.Bookmarks(profileName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	yamlBookmarks, err := yaml.Marshal(bookmarks)
	if err != nil {
		return NewTextResult("", err), nil
	}
	return NewTextResult(fmt.Sprintf("The following bookmarks (YAML format) were found:\n%s", string(yamlBookmarks)), nil), nil
}
