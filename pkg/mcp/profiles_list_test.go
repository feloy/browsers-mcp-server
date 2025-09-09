package mcp

import (
	"context"
	"testing"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/test"
	"github.com/feloy/browsers-mcp-server/pkg/config"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestListProfiles(t *testing.T) {

	var browser1 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser1",
		Available: true,
		Profiles:  []string{"profile1a"},
		Bookmarks: []api.BookMark{},
	})

	var browser2 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser2",
		Available: false,
		Profiles:  []string{"profile2"},
		Bookmarks: []api.BookMark{},
	})

	var browser3 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser3",
		Available: true,
		Profiles:  []string{"profile3a", "profile3b"},
		Bookmarks: []api.BookMark{},
	})

	for _, tt := range []struct {
		name       string
		browsers   []*test.Browser
		parameters map[string]interface{}
		expected   string
	}{
		{
			name:     "one available browser with one profile, no parameter passed",
			browsers: []*test.Browser{browser1, browser2},
			expected: "profile1a",
		},
		{
			name:     "one available browser with two profiles, no parameter passed",
			browsers: []*test.Browser{browser3, browser2},
			expected: "profile3a, profile3b",
		},
		{
			name:     "two available browsers, no parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			expected: "failed to get list of profiles, multiple browsers found, please specify the browser",
		},
		{
			name:     "two available browsers, browser parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser1",
			},
			expected: "profile1a",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			browsers.Clear()
			for _, browser := range tt.browsers {
				browsers.Register(browser)
			}
			server, err := NewServer(Configuration{
				Profile: &FullProfile{},
				StaticConfig: &config.StaticConfig{
					EnabledTools: []string{"list_profiles"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}
			ctr := mcp.CallToolRequest{}
			ctr.Params.Arguments = tt.parameters
			result, err := server.listProfiles(context.Background(), ctr)
			if err != nil {
				t.Fatalf("Failed to list browsers: %v", err)
			}
			if len(result.Content) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(result.Content))
			}
			text := result.Content[0].(mcp.TextContent).Text
			if text != tt.expected {
				t.Fatalf("Expected %s, got %s", tt.expected, text)
			}
		})
	}
}
