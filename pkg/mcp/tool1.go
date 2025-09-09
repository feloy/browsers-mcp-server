package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/feloy/mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
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
		{
			Tool: mcp.NewTool("list_bookmarks",
				mcp.WithDescription("List the available bookmarks for a browser's profile"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the bookmarks for"),
					mcp.Required(),
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

func (s *Server) listBookmarks(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var browserName string
	ok := false
	if browserName, ok = ctr.GetArguments()["browser"].(string); !ok {
		return NewTextResult("", fmt.Errorf("failed to get list of bookmarks, missing argument browser")), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	var profileName string
	ok = false
	if profileName, ok = ctr.GetArguments()["profile"].(string); !ok {
		profiles, err := browser.Profiles()
		if err != nil {
			return NewTextResult("", err), nil
		}
		if len(profiles) == 1 {
			profileName = profiles[0]
		} else {
			return NewTextResult("", fmt.Errorf("failed to get list of bookmarks, multiple profiles found, please specify the profile")), nil
		}
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
