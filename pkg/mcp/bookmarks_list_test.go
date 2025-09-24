package mcp

import (
	"context"
	"slices"
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

func TestListBookmarks(t *testing.T) {

	var browser1 = test.NewBrowser(test.NewBrowserOptions{
		Name:      "browser1",
		Available: true,
		Profiles:  []string{"profile1a"},
		Bookmarks: []api.BookMark{
			{
				Name: "bookmark1a", URL: "https://www.bookmark1a.com", Folder: []string{"folder1a"},
				DateAdded:       globaltest.Must(time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")),
				DateLastVisited: globaltest.Must(time.Parse(time.RFC3339, "2021-04-01T00:00:00Z")),
			},
		},
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
		Bookmarks: []api.BookMark{
			{
				Name: "bookmark3a", URL: "https://www.bookmark3a.com", Folder: []string{"folder3a"},
				DateAdded:       globaltest.Must(time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")),
				DateModified:    globaltest.Must(time.Parse(time.RFC3339, "2023-02-01T00:00:00Z")),
				DateLastVisited: globaltest.Must(time.Parse(time.RFC3339, "2023-04-01T00:00:00Z")),
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
		expected_input_properties_required     [][]bool
		expected_input_properties_descriptions [][]string

		toolName   string
		parameters map[string]interface{}
		expected   string
	}{
		{
			name:                  "one available browser with one profile",
			browsers:              []*test.Browser{browser1, browser2},
			expected_tools_count:  1,
			expected_names:        []string{"list_bookmarks"},
			expected_descriptions: []string{"List the available bookmarks in the browser"},
			expected_input_properties: [][]string{
				{},
			},
			expected_input_properties_required: [][]bool{
				{},
			},
			expected_input_properties_descriptions: [][]string{
				{},
			},
			toolName:   "list_bookmarks",
			parameters: map[string]interface{}{},
			expected: `The following bookmarks (YAML format) were found:
- name: bookmark1a
  url: https://www.bookmark1a.com
  folder:
    - folder1a
  date_added: 2021-01-01T00:00:00Z
  date_last_visited: 2021-04-01T00:00:00Z
`,
		},
		{
			name:                  "one available browser with two profiles",
			browsers:              []*test.Browser{browser2, browser3},
			expected_tools_count:  1,
			expected_names:        []string{"list_bookmarks"},
			expected_descriptions: []string{"List the available bookmarks in the browser"},
			expected_input_properties: [][]string{
				{"profile"},
			},
			expected_input_properties_required: [][]bool{
				{true},
			},
			expected_input_properties_descriptions: [][]string{
				{"The browser's profile to list the bookmarks for"},
			},
			toolName: "list_bookmarks",
			parameters: map[string]interface{}{
				"profile": "profile3a",
			},
			expected: `The following bookmarks (YAML format) were found:
- name: bookmark3a
  url: https://www.bookmark3a.com
  folder:
    - folder3a
  date_added: 2023-01-01T00:00:00Z
  date_modified: 2023-02-01T00:00:00Z
  date_last_visited: 2023-04-01T00:00:00Z
`,
		},
		{
			name:                 "two available browsers with one or two profiles",
			browsers:             []*test.Browser{browser1, browser2, browser3},
			expected_tools_count: 1,
			expected_names:       []string{"list_bookmarks"},
			expected_descriptions: []string{
				"List the available bookmarks in the browser",
			},
			expected_input_properties: [][]string{
				{"profile"},
			},
			expected_input_properties_required: [][]bool{
				{true},
			},
			expected_input_properties_descriptions: [][]string{
				{"The browser's profile to list the bookmarks for"},
			},
			toolName: "list_bookmarks",
			parameters: map[string]interface{}{
				"profile": "profile3a on browser3",
			},
			expected: `The following bookmarks (YAML format) were found:
- name: bookmark3a
  url: https://www.bookmark3a.com
  folder:
    - folder3a
  date_added: 2023-01-01T00:00:00Z
  date_modified: 2023-02-01T00:00:00Z
  date_last_visited: 2023-04-01T00:00:00Z
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
					EnabledTools: []string{"list_bookmarks"},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create server: %v", err)
			}

			// test API
			tools := srv.initBookmarksList()
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

					required := slices.Contains(tool.Tool.InputSchema.Required, property)
					if required != tt.expected_input_properties_required[i][j] {
						t.Fatalf("expected property required %v, got %v", tt.expected_input_properties_required[i][j], required)
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
