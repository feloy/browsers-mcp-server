package mcp

import (
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *Server) getBrowserName(ctr mcp.CallToolRequest, info string) (string, error) {
	var browserName string
	ok := false
	if browserName, ok = ctr.GetArguments()["browser"].(string); !ok {
		browsers := browsers.GetBrowsers()
		if len(browsers) == 1 {
			browserName = browsers[0].Name()
		} else {
			return "", fmt.Errorf("failed to get %s, multiple browsers found, please specify the browser", info)
		}
	}
	return browserName, nil
}

func (s *Server) getProfileName(browser api.Browser, ctr mcp.CallToolRequest, info string) (string, error) {
	var profileName string
	ok := false
	if profileName, ok = ctr.GetArguments()["profile"].(string); !ok {
		profiles, err := browser.Profiles()
		if err != nil {
			return "", err
		}
		if len(profiles) == 1 {
			profileName = profiles[0]
		} else {
			return "", fmt.Errorf("failed to get %s, multiple profiles found, please specify the profile", info)
		}
	}
	return profileName, nil
}
