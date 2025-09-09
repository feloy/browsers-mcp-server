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

func TestListBrowsers(t *testing.T) {

	var browser1 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser1",
		Available: true,
		Profiles:  []string{"profile1"},
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
		Profiles:  []string{"profile3"},
		Bookmarks: []api.BookMark{},
	})

	for _, tt := range []struct {
		name     string
		browsers []*test.Browser
		expected string
	}{
		{
			name:     "one available browser",
			browsers: []*test.Browser{browser1, browser2},
			expected: "browser1",
		},
		{
			name:     "two available browsers",
			browsers: []*test.Browser{browser1, browser3},
			expected: "browser1, browser3",
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
					EnabledTools: []string{"list_browsers"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}
			result, err := server.listBrowsers(context.Background(), mcp.CallToolRequest{})
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
