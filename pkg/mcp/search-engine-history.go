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
					mcp.Description("List the search engine queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight"),
				),
				mcp.WithNumber(
					"limit",
					mcp.Description("The maximum number of search engine queries to list, default is 10"),
					mcp.DefaultNumber(10),
				),
			),
			Handler: s.listSearchEngineQueries,
		}, {
			Tool: mcp.NewTool("list_visited_pages_from_search_engine_query",
				mcp.WithDescription("list visited pages from a search engine query"),
				mcp.WithString(
					"browser",
					mcp.Description("The browser to list the visited pages for"),
				),
				mcp.WithString(
					"profile",
					mcp.Description("The browser's profile to list the visited pages for, required if there are multiple profiles"),
				),
				mcp.WithString(
					"query",
					mcp.Description("The query string to list the visited pages for"),
					mcp.Required(),
				),
				mcp.WithString(
					"start_time",
					mcp.Description("List the visited pages for queries from this time (YYYY-MM-DD HH:MM:SS), default is today at midnight"),
				),
			),
			Handler: s.listVisitedPagesFromSearchEngineQuery,
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
		t, err := time.Parse(time.DateTime, startTimeStr)
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

func (s *Server) listVisitedPagesFromSearchEngineQuery(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	browserName, err := s.getBrowserName(ctr, "visited pages from search engine query")
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	profileName, err := s.getProfileName(browser, ctr, "visited pages from search engine query")
	if err != nil {
		return NewTextResult("", err), nil
	}

	var startTime *time.Time
	if startTimeStr, ok := ctr.GetArguments()["start_time"].(string); ok {
		t, err := time.Parse(time.DateTime, startTimeStr)
		if err != nil {
			return NewTextResult("", err), nil
		}
		startTime = &t
	}

	query, ok := ctr.GetArguments()["query"].(string)
	if !ok {
		return NewTextResult("", fmt.Errorf("query is required")), nil
	}

	visitedPages, err := browser.ListVisitedPagesFromSearchEngineQuery(profileName, api.ListVisitedPagesFromSearchEngineQueryOptions{StartTime: startTime, Query: query})
	if err != nil {
		return NewTextResult("", err), nil
	}

	yamlVisitedPages, err := yaml.Marshal(visitedPages)
	if err != nil {
		return NewTextResult("", err), nil
	}

	return NewTextResult(fmt.Sprintf("The following visited pages (YAML format) were found:\n%s", string(yamlVisitedPages)), nil), nil
}
