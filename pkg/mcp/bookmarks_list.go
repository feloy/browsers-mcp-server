package mcp

import (
	"context"
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

type ListBookmarksResult struct {
	Browser   string         `json:"browser_name"`
	Profile   string         `json:"profile_name"`
	Bookmarks []api.BookMark `json:"bookmarks"`
}

func (s *Server) initBookmarksList() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_bookmarks",
				mcp.WithDescription("List the available bookmarks for a browser's profile. For each bookmark, the Name is the information to display to the user. The URL is the information to use to navigate to this bookmark. The folder indicates in which folder the bookmark is located in the bookmarks categorization."),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the bookmarks for. The list of browsers is returned by the `list_browsers` tool."),
				),
				mcp.WithString(
					"profile",
					mcp.Description("The browser's profile to list the bookmarks for, required if there are multiple profiles. The list of profiles is returned by the `list_profiles` tool."),
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
	return NewStructuredResult(fmt.Sprintf("The following bookmarks (YAML format) were found:\n%s", string(yamlBookmarks)), ListBookmarksResult{
		Browser:   browserName,
		Profile:   profileName,
		Bookmarks: bookmarks,
	}, nil), nil
}
