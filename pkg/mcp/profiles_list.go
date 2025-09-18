package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListProfilesResult struct {
	Browser  string   `json:"browser_name"`
	Profiles []string `json:"profiles"`
}

func (s *Server) initProfilesList() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_profiles",
				mcp.WithDescription("List the available profiles for a browser"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the profiles for. The list of browsers is returned by the `list_browsers` tool."),
				),
			),
			Handler: s.listProfiles,
		},
	}
	return tools
}

func (s *Server) listProfiles(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var browserName string
	ok := false
	if browserName, ok = ctr.GetArguments()["browser"].(string); !ok {
		browsers := browsers.GetBrowsers()
		if len(browsers) == 1 {
			browserName = browsers[0].Name()
		} else {
			return NewTextResult("", fmt.Errorf("failed to get list of profiles, multiple browsers found, please specify the browser")), nil
		}
	}

	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}
	profiles, err := browser.Profiles()
	if err != nil {
		return NewTextResult("", err), nil
	}
	return NewStructuredResult(strings.Join(profiles, ", "), ListProfilesResult{
		Browser:  browserName,
		Profiles: profiles,
	}, nil), nil
}
