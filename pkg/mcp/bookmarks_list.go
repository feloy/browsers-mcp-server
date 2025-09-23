package mcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initBookmarksList() []server.ServerTool {
	browsers := browsers.GetBrowsers()
	if len(browsers) == 1 {
		options := []mcp.ToolOption{
			mcp.WithDescription("List the available bookmarks in the browser"),
		}
		profiles, err := browsers[0].Profiles()
		if err != nil {
			return nil
		}
		if len(profiles) > 1 {
			options = append(options,
				mcp.WithString(
					"profile",
					mcp.Required(),
					mcp.Enum(profiles...),
					mcp.Description("The browser's profile to list the bookmarks for"),
				))
		}
		return []server.ServerTool{
			{
				Tool:    mcp.NewTool("list_bookmarks", options...),
				Handler: s.listBookmarksForBrowser(nil),
			},
		}
	} else {
		var tools []server.ServerTool
		for _, browser := range browsers {
			options := []mcp.ToolOption{
				mcp.WithDescription(fmt.Sprintf("List the available bookmarks in browser %s", browser.Name())),
			}
			profiles, err := browser.Profiles()
			if err != nil {
				return nil
			}
			if len(profiles) > 1 {
				options = append(options,
					mcp.WithString(
						"profile",
						mcp.Required(),
						mcp.Enum(profiles...),
						mcp.Description("The browser's profile to list the bookmarks for"),
					),
				)
			}
			browserName := browser.Name()
			tools = append(tools, server.ServerTool{
				Tool:    mcp.NewTool(fmt.Sprintf("list_bookmarks_%s", browser.Name()), options...),
				Handler: s.listBookmarksForBrowser(&browserName),
			})
		}
		return tools
	}
}

func (s *Server) listBookmarksForBrowser(browserName *string) func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if browserName == nil {
			browsersList := browsers.GetBrowsers()
			if len(browsersList) != 1 {
				return nil, errors.New("more than one browser found, this is not expected")
			}
			name := browsersList[0].Name()
			browserName = &name
		}
		browser, err := browsers.GetBrowserByName(*browserName)
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
}
