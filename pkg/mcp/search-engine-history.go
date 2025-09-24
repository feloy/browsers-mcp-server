package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initSearchEngineQueries() []server.ServerTool {
	tools := []server.ServerTool{}
	listSearchEngineTool, err := s.getListSearchEngineQueries()
	if err == nil {
		tools = append(tools, listSearchEngineTool...)
	}
	listVisitedPagesFromSearchEngineQueryTool, err := s.getListVisitedPagesFromSearchEngineQuery()
	if err == nil {
		tools = append(tools, listVisitedPagesFromSearchEngineQueryTool...)
	}

	var toolsNames []string
	for _, tool := range tools {
		toolsNames = append(toolsNames, tool.Tool.Name)
	}
	log.Info(fmt.Sprintf("registering %d tools in initSearchEngineQueries: %s", len(tools), strings.Join(toolsNames, ", ")))
	return tools
}

func (s *Server) getListSearchEngineQueries() ([]server.ServerTool, error) {
	options := []mcp.ToolOption{
		mcp.WithDescription("list queries in search engines"),
	}

	browserProfiles := BrowsersProfiles{}
	browserProfiles.Populate(browsers.GetBrowsers())
	profilesEnum := browserProfiles.FlatList()

	if len(profilesEnum) > 0 {
		options = append(options,
			mcp.WithString(
				"profile",
				mcp.Required(),
				mcp.Enum(profilesEnum...),
				mcp.Description("The browser's profile to list the search engine queries for"),
			))
	}

	options = append(
		options,
		mcp.WithString(
			"day",
			mcp.Description("List the search engine queries done on this day (YYYY-MM-DD), default is today"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("The maximum number of search engine queries to list, default is 10"),
			mcp.DefaultNumber(10),
		),
	)
	return []server.ServerTool{
		{
			Tool:    mcp.NewTool("list_search_engine_queries", options...),
			Handler: s.listSearchEnginesQueries,
		},
	}, nil

}

func (s *Server) getListVisitedPagesFromSearchEngineQuery() ([]server.ServerTool, error) {
	options := []mcp.ToolOption{
		mcp.WithDescription("list the pages visited after doing a specific query in a search engine"),
	}

	browserProfiles := BrowsersProfiles{}
	browserProfiles.Populate(browsers.GetBrowsers())
	profilesEnum := browserProfiles.FlatList()

	if len(profilesEnum) > 0 {
		options = append(options,
			mcp.WithString(
				"profile",
				mcp.Required(),
				mcp.Enum(profilesEnum...),
				mcp.Description("The browser's profile to list the visited pages for"),
			))
	}

	options = append(
		options,
		mcp.WithString(
			"query",
			mcp.Description("The query string to list the visited pages for"),
			mcp.Required(),
		),
		mcp.WithString(
			"day",
			mcp.Description("List the visited pages for queries done on this day (YYYY-MM-DD), default is today"),
		),
	)
	return []server.ServerTool{
		{
			Tool:    mcp.NewTool("list_visited_pages_from_search_engine_query", options...),
			Handler: s.listVisitedPagesFromSearchEngineQuery,
		},
	}, nil

}

func (s *Server) listSearchEnginesQueries(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	profileParam, _ := ctr.GetArguments()["profile"].(string)
	browserName, profileName, err := GetBrowserAndProfileFromValue(profileParam, browsers.GetBrowsers())
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	var startTime time.Time
	var endTime time.Time
	if startDayStr, ok := ctr.GetArguments()["day"].(string); ok {
		t, err := time.Parse(time.DateOnly, startDayStr)
		if err != nil {
			return NewTextResult("", err), nil
		}
		startTime = t
		endTime = t.AddDate(0, 0, 1)
	} else {
		startTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
		endTime = startTime.AddDate(0, 0, 1)
	}

	var limit int
	if limitFloat, ok := ctr.GetArguments()["limit"].(float64); ok {
		limit = int(limitFloat)
	} else {
		// should be set because of `mcp.DefaultNumber(10)` but this is not the case
		limit = 10
	}

	searchEngineQueries, err := browser.SearchEngineQueries(profileName, api.SearchEngineOptions{StartTime: startTime, EndTime: endTime, Limit: limit})
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
	profileParam, _ := ctr.GetArguments()["profile"].(string)
	browserName, profileName, err := GetBrowserAndProfileFromValue(profileParam, browsers.GetBrowsers())
	if err != nil {
		return NewTextResult("", err), nil
	}
	browser, err := browsers.GetBrowserByName(browserName)
	if err != nil {
		return NewTextResult("", err), nil
	}

	var startTime time.Time
	var endTime time.Time
	if startDayStr, ok := ctr.GetArguments()["day"].(string); ok {
		t, err := time.Parse(time.DateOnly, startDayStr)
		if err != nil {
			return NewTextResult("", err), nil
		}
		startTime = t
		endTime = t.AddDate(0, 0, 1)
	} else {
		startTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
		endTime = startTime.AddDate(0, 0, 1)
	}

	query, ok := ctr.GetArguments()["query"].(string)
	if !ok {
		return NewTextResult("", fmt.Errorf("query is required")), nil
	}

	visitedPages, err := browser.ListVisitedPagesFromSearchEngineQuery(profileName, api.ListVisitedPagesFromSearchEngineQueryOptions{StartTime: startTime, EndTime: endTime, Query: query})
	if err != nil {
		return NewTextResult("", err), nil
	}

	yamlVisitedPages, err := yaml.Marshal(visitedPages)
	if err != nil {
		return NewTextResult("", err), nil
	}

	return NewTextResult(fmt.Sprintf("The following visited pages (YAML format) were found:\n%s", string(yamlVisitedPages)), nil), nil
}
