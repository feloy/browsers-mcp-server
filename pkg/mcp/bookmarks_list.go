package mcp

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initBookmarksList() []server.ServerTool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List the available bookmarks in the browser"),
	}

	browserProfiles := BrowsersProfiles{}
	browserProfiles.Populate(browsers.GetBrowsers())
	profilesEnum := browserProfiles.FlatList()
	log.Debug("bookmarks list", "profilesEnum", profilesEnum)

	if len(profilesEnum) > 0 {
		options = append(options,
			mcp.WithString(
				"profile",
				mcp.Required(),
				mcp.Enum(profilesEnum...),
				mcp.Description("The browser's profile to list the bookmarks for"),
			))
	}
	return []server.ServerTool{
		{
			Tool:    mcp.NewTool("list_bookmarks", options...),
			Handler: s.listBookmarks,
		},
	}

}

func (s *Server) listBookmarks(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profileParam, _ := ctr.GetArguments()["profile"].(string)
	browserName, profileName, err := GetBrowserAndProfileFromValue(profileParam, browsers.GetBrowsers())
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
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
