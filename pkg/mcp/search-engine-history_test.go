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
	"github.com/mark3labs/mcp-go/server"
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
		name                                   string
		browsers                               []*test.Browser
		expected_tools_count                   int
		expected_names                         []string
		expected_descriptions                  []string
		expected_input_properties              [][]string
		expected_input_properties_descriptions [][]string

		toolName   string
		parameters map[string]interface{}
		expected   string
	}{
		{
			name:                 "one available browser with one profile",
			browsers:             []*test.Browser{browser1, browser2},
			expected_tools_count: 2,
			expected_names: []string{
				"list_search_engine_queries",
				"list_visited_pages_from_search_engine_query",
			},
			expected_descriptions: []string{
				"list queries in search engines",
				"list the pages visited after doing a specific query in a search engine",
			},
			expected_input_properties: [][]string{
				{"day", "limit"},
				{"query", "day"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"List the search engine queries done on this day (YYYY-MM-DD), default is today",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The query string to list the visited pages for",
					"List the visited pages for queries done on this day (YYYY-MM-DD), default is today",
				},
			},
			toolName:   "list_search_engine_queries",
			parameters: map[string]interface{}{},
			expected: `The following search queries (YAML format) were found:
- query: where is charly
  date: 2021-01-01T00:00:00Z
  search_engine: Google
`,
		},

		{
			name:                 "one available browser with several profiles",
			browsers:             []*test.Browser{browser2, browser3},
			expected_tools_count: 2,
			expected_names: []string{
				"list_search_engine_queries",
				"list_visited_pages_from_search_engine_query",
			},
			expected_descriptions: []string{
				"list queries in search engines",
				"list the pages visited after doing a specific query in a search engine",
			},
			expected_input_properties: [][]string{
				{"profile", "day", "limit"},
				{"profile", "query", "day"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"The browser's profile to list the search engine queries for, possible values are profile3a, profile3b",
					"List the search engine queries done on this day (YYYY-MM-DD), default is today",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The browser's profile to list the visited pages for, possible values are profile3a, profile3b",
					"The query string to list the visited pages for",
					"List the visited pages for queries done on this day (YYYY-MM-DD), default is today",
				},
			},
			toolName: "list_search_engine_queries",
			parameters: map[string]interface{}{
				"profile": "profile3a",
			},
			expected: `The following search queries (YAML format) were found:
- query: what is it
  date: 2023-01-01T00:00:00Z
  search_engine: Google
`,
		},
		{
			name:                 "two available browsers with one or several profiles",
			browsers:             []*test.Browser{browser1, browser2, browser3},
			expected_tools_count: 4,
			expected_names: []string{
				"list_search_engine_queries_browser1",
				"list_search_engine_queries_browser3",
				"list_visited_pages_from_search_engine_query_browser1",
				"list_visited_pages_from_search_engine_query_browser3",
			},
			expected_descriptions: []string{
				"list queries in search engines in browser browser1",
				"list queries in search engines in browser browser3",
				"list the pages visited after doing a specific query in a search engine in browser browser1",
				"list the pages visited after doing a specific query in a search engine in browser browser3",
			},
			expected_input_properties: [][]string{
				{"day", "limit"},
				{"profile", "day", "limit"},
				{"query", "day"},
				{"profile", "query", "day"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"List the search engine queries done on this day (YYYY-MM-DD), default is today",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The browser3's profile to list the search engine queries for, possible values are profile3a, profile3b",
					"List the search engine queries done on this day (YYYY-MM-DD), default is today",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The query string to list the visited pages for",
					"List the visited pages for queries done on this day (YYYY-MM-DD), default is today",
				},
				{
					"The browser3's profile to list the visited pages for, possible values are profile3a, profile3b",
					"The query string to list the visited pages for",
					"List the visited pages for queries done on this day (YYYY-MM-DD), default is today",
				},
			},
			toolName: "list_search_engine_queries_browser3",
			parameters: map[string]interface{}{
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
			srv, err := NewServer(Configuration{
				Profile: &FullProfile{},
				StaticConfig: &config.StaticConfig{
					EnabledTools: []string{"list_search_engine_queries"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}
			tools := srv.initSearchEngineQueries()
			if len(tools) != tt.expected_tools_count {
				t.Fatalf("Expected %d tools, got %d", tt.expected_tools_count, len(tools))
			}
			// Check API
			for i, tool := range tools {
				if tool.Tool.Name != tt.expected_names[i] {
					t.Fatalf("Expected tool name #%d to be %s, but is %s", i, tt.expected_names[i], tool.Tool.Name)
				}
				if tool.Tool.Description != tt.expected_descriptions[i] {
					t.Fatalf("Expected tool description #%d to be %s, but is %s", i, tt.expected_descriptions[i], tool.Tool.Description)
				}
				if len(tool.Tool.InputSchema.Properties) != len(tt.expected_input_properties[i]) {
					t.Fatalf("Expected input properties count #%d to be %d, but is %d", i, len(tt.expected_input_properties[i]), len(tool.Tool.InputSchema.Properties))
				}
				for j, property := range tt.expected_input_properties[i] {
					var foundProperty any
					var found bool
					if foundProperty, found = tool.Tool.InputSchema.Properties[property]; !found {
						t.Fatalf("expected property %s not found", property)
					}
					var option map[string]any
					var ok bool
					if option, ok = foundProperty.(map[string]any); !ok {
						t.Fatalf("cast error")
					}
					var description string
					if description, ok = option["description"].(string); !ok {
						t.Fatalf("description cast error")
					}
					if description != tt.expected_input_properties_descriptions[i][j] {
						t.Fatalf("expected property description %q, got %q", tt.expected_input_properties_descriptions[i][j], description)
					}
				}
			}

			// Test call
			var tool server.ServerTool
			var found = false
			for _, tool = range tools {
				if tool.Tool.Name == tt.toolName {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("tool %s not found", tt.toolName)
			}
			ctr := mcp.CallToolRequest{}
			ctr.Params.Arguments = tt.parameters
			var result *mcp.CallToolResult
			result, err = tool.Handler(context.Background(), ctr)
			if err != nil {
				t.Fatalf("Failed to call tool: %v", err)
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
