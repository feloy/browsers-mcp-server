package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/test"
	"github.com/feloy/browsers-mcp-server/pkg/config"
	globaltest "github.com/feloy/browsers-mcp-server/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestListSearchEngineHistory(t *testing.T) {

	var browser1 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser1",
		Available: true,
		Profiles:  []string{"profile1a"},
		SearchEngineQueries: []api.SearchEngineQuery{
			{
				Query:        "where is charly",
				Date:         globaltest.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")),
				SearchEngine: "Google",
			},
		},
	})

	var browser2 = test.NewBrowser(test.NewBrowserOptions{
		Name:                "browser2",
		Available:           false,
		Profiles:            []string{"profile2"},
		SearchEngineQueries: []api.SearchEngineQuery{},
	})

	var browser3 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser3",
		Available: true,
		Profiles:  []string{"profile3a", "profile3b"},
		SearchEngineQueries: []api.SearchEngineQuery{
			{
				Query:        "what is it",
				Date:         globaltest.Must(time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")),
				SearchEngine: "Google",
			},
		},
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
			expected: `The following search queries (YAML format) were found:
- query: where is charly
  date: 2021-01-01T00:00:00Z
  search_engine: Google
`,
		},
		{
			name:     "one available browser with two profiles, no parameter passed",
			browsers: []*test.Browser{browser3, browser2},
			expected: "failed to get search engine queries, multiple profiles found, please specify the profile",
		},
		{
			name:     "two available browsers, no parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			expected: "failed to get search engine queries, multiple browsers found, please specify the browser",
		},
		{
			name:     "two available browsers, browser with 1 profile passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser1",
			},
			expected: `The following search queries (YAML format) were found:
- query: where is charly
  date: 2021-01-01T00:00:00Z
  search_engine: Google
`,
		},
		{
			name:     "two available browsers, browser with 2 profiles passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser3",
			},
			expected: "failed to get search engine queries, multiple profiles found, please specify the profile",
		},
		{
			name:     "two available browsers, browser with 2 profiles passed, profile parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser3",
				"profile": "profile3a",
			},
			expected: `The following search queries (YAML format) were found:
- query: what is it
  date: 2023-01-01T00:00:00Z
  search_engine: Google
`,
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
					EnabledTools: []string{"list_bookmarks"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}
			ctr := mcp.CallToolRequest{}
			ctr.Params.Arguments = tt.parameters
			result, err := server.listSearchEngineQueries(context.Background(), ctr)
			if err != nil {
				t.Fatalf("Failed to list search engine queries: %v", err)
			}
			if len(result.Content) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(result.Content))
			}
			text := result.Content[0].(mcp.TextContent).Text
			if text != tt.expected {
				t.Fatalf("Content differs:\n%s", cmp.Diff(tt.expected, text))
			}
		})
	}
}
