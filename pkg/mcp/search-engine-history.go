package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initSearchEngineQueries() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("list_search_engine_queries",
				mcp.WithDescription("list queries in search engines"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the search engine queries for"),
				),
				mcp.WithString(
					"profile",
					mcp.Description("The browser's profile to list the search engine queries for, required if there are multiple profiles"),
				),
				mcp.WithString(
					"start_time",
					mcp.Description("List the search engine queries from this time (YYYY-MM-DD:HH-MM-SS), default is today at midnight"),
				),
				mcp.WithNumber(
					"limit",
					mcp.Description("The maximum number of search engine queries to list, default is 10"),
					mcp.DefaultNumber(10),
				),
			),
			Handler: s.listSearchEngineQueries,
		},
	}
	return tools
}

func (s *Server) listSearchEngineQueries(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	browserName, err := s.getBrowserName(ctr, "search engine queries")
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	profileName, err := s.getProfileName(browser, ctr, "search engine queries")
	if err != nil {
		return NewTextResult("", err), nil
	}

	var startTime *time.Time
	if startTimeStr, ok := ctr.GetArguments()["start_time"].(string); ok {
		t, err := time.Parse("2006-01-02:15-04-05", startTimeStr)
		if err != nil {
			return NewTextResult("", err), nil
		}
		startTime = &t
	}

	var limit int
	if limitFloat, ok := ctr.GetArguments()["limit"].(float64); ok {
		limit = int(limitFloat)
	}

	searchEngineQueries, err := browser.SearchEngineQueries(profileName, api.SearchEngineOptions{StartTime: startTime, Limit: limit})
	if err != nil {
		return NewTextResult("", err), nil
	}

	yamlSearchEngineQueries, err := yaml.Marshal(searchEngineQueries)
	if err != nil {
		return NewTextResult("", err), nil
	}

	return NewTextResult(fmt.Sprintf("The following search queries (YAML format) were found:\n%s", string(yamlSearchEngineQueries)), nil), nil
}
