package mcp

import (
	"testing"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/test"
	"github.com/feloy/browsers-mcp-server/pkg/config"
	globaltest "github.com/feloy/browsers-mcp-server/pkg/test"
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
	_ = browser3

	for _, tt := range []struct {
		name                                   string
		browsers                               []*test.Browser
		parameters                             map[string]interface{}
		expected_tools_count                   int
		expected_names                         []string
		expected_descriptions                  []string
		expected_input_properties              [][]string
		expected_input_properties_descriptions [][]string
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
				{"start_time", "limit"},
				{"query", "start_time"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"List the search engine queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The query string to list the visited pages for",
					"List the visited pages for queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
				},
			}},
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
				{"profile", "start_time", "limit"},
				{"profile", "query", "start_time"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"The browser's profile to list the search engine queries for, possible values are profile3a, profile3b",
					"List the search engine queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The browser's profile to list the visited pages for, possible values are profile3a, profile3b",
					"The query string to list the visited pages for",
					"List the visited pages for queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
				},
			},
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
				{"start_time", "limit"},
				{"profile", "start_time", "limit"},
				{"query", "start_time"},
				{"profile", "query", "start_time"},
			},
			expected_input_properties_descriptions: [][]string{
				{
					"List the search engine queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The browser3's profile to list the search engine queries for, possible values are profile3a, profile3b",
					"List the search engine queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
					"The maximum number of search engine queries to list, default is 10",
				},
				{
					"The query string to list the visited pages for",
					"List the visited pages for queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
				},
				{
					"The browser3's profile to list the visited pages for, possible values are profile3a, profile3b",
					"The query string to list the visited pages for",
					"List the visited pages for queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight",
				},
			},
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
					EnabledTools: []string{"list_search_engine_queries"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}
			tools := server.initSearchEngineQueries()
			if len(tools) != tt.expected_tools_count {
				t.Fatalf("Expected %d tools, got %d", tt.expected_tools_count, len(tools))
			}
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
		})
	}
}
