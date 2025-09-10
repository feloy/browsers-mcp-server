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
		name       string
		browsers   []*test.Browser
		parameters map[string]interface{}
		expected   string
	}{
		{
			name:     "one available browser with one profile, no parameter passed",
			browsers: []*test.Browser{browser1, browser2},
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
			name:     "one available browser with two profiles, no parameter passed",
			browsers: []*test.Browser{browser3, browser2},
			expected: "failed to get list of bookmarks, multiple profiles found, please specify the profile",
		},
		{
			name:     "two available browsers, no parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			expected: "failed to get list of bookmarks, multiple browsers found, please specify the browser",
		},
		{
			name:     "two available browsers, browser with 1 profile passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser1",
			},
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
			name:     "two available browsers, browser with 2 profiles passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser3",
			},
			expected: "failed to get list of bookmarks, multiple profiles found, please specify the profile",
		},
		{
			name:     "two available browsers, browser with 2 profiles passed, profile parameter passed",
			browsers: []*test.Browser{browser1, browser3},
			parameters: map[string]interface{}{
				"browser": "browser3",
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
			result, err := server.listBookmarks(context.Background(), ctr)
			if err != nil {
				t.Fatalf("Failed to list bookmarks: %v", err)
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
