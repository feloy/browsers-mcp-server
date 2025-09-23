package mcp

import (
	"context"
	"errors"
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
	browsers := browsers.GetBrowsers()
	if len(browsers) == 1 {
		profiles, err := browsers[0].Profiles()
		if err != nil {
			return nil, err
		}
		options := []mcp.ToolOption{}
		options = append(options, mcp.WithDescription("list queries in search engines"))
		if len(profiles) > 1 {
			options = append(options, mcp.WithString("profile",
				mcp.Required(),
				mcp.Enum(profiles...),
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
				Handler: s.listSearchEnginesQueriesByBrowser(nil),
			},
		}, nil
	} else {
		var tools []server.ServerTool
		for _, browser := range browsers {
			browserName := browser.Name()
			profiles, err := browser.Profiles()
			if err != nil {
				return nil, err
			}
			options := []mcp.ToolOption{}
			options = append(options, mcp.WithDescription(fmt.Sprintf("list queries in search engines in browser %s", browser.Name())))
			if len(profiles) > 1 {
				options = append(options, mcp.WithString("profile",
					mcp.Required(),
					mcp.Enum(profiles...),
					mcp.Description(fmt.Sprintf("The %s's profile to list the search engine queries for", browser.Name())),
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
			tools = append(tools, server.ServerTool{
				Tool:    mcp.NewTool(fmt.Sprintf("list_search_engine_queries_%s", browser.Name()), options...),
				Handler: s.listSearchEnginesQueriesByBrowser(&browserName),
			})
		}
		return tools, nil
	}
}

func (s *Server) getListVisitedPagesFromSearchEngineQuery() ([]server.ServerTool, error) {
	browsers := browsers.GetBrowsers()
	if len(browsers) == 1 {
		profiles, err := browsers[0].Profiles()
		if err != nil {
			return nil, err
		}
		options := []mcp.ToolOption{}
		options = append(options, mcp.WithDescription("list the pages visited after doing a specific query in a search engine"))
		if len(profiles) > 1 {
			options = append(options, mcp.WithString("profile",
				mcp.Required(),
				mcp.Enum(profiles...),
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
				Handler: s.listVisitedPagesFromSearchEngineQueryByBrowser(nil),
			},
		}, nil
	} else {
		var tools []server.ServerTool
		for _, browser := range browsers {
			browserName := browser.Name()
			profiles, err := browser.Profiles()
			if err != nil {
				return nil, err
			}
			options := []mcp.ToolOption{}
			options = append(options, mcp.WithDescription(fmt.Sprintf("list the pages visited after doing a specific query in a search engine in browser %s", browser.Name())))
			if len(profiles) > 1 {
				options = append(options, mcp.WithString("profile",
					mcp.Required(),
					mcp.Enum(profiles...),
					mcp.Description(fmt.Sprintf("The %s's profile to list the visited pages for", browser.Name())),
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
			tools = append(tools, server.ServerTool{
				Tool:    mcp.NewTool(fmt.Sprintf("list_visited_pages_from_search_engine_query_%s", browser.Name()), options...),
				Handler: s.listVisitedPagesFromSearchEngineQueryByBrowser(&browserName),
			})
		}
		return tools, nil
	}
}

func (s *Server) listSearchEnginesQueriesByBrowser(browserName *string) func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if browserName == nil {
			browsersList := browsers.GetBrowsers()
			if len(browsersList) != 1 {
				return nil, errors.New("more than one browser found, this is not expected")
			}
			name := browsersList[0].Name()
			browserName = &name
		}
		browser, err := browsers.GetBrowserByName(*browserName)
		if err != nil {
			return NewTextResult("", err), nil
		}

		profileName, err := s.getProfileName(browser, ctr, "search engine queries")
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
}

func (s *Server) listVisitedPagesFromSearchEngineQueryByBrowser(browserName *string) func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	return func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if browserName == nil {
			browsersList := browsers.GetBrowsers()
			if len(browsersList) != 1 {
				return nil, errors.New("more than one browser found, this is not expected")
			}
			name := browsersList[0].Name()
			browserName = &name
		}
		browser, err := browsers.GetBrowserByName(*browserName)
		if err != nil {
			return NewTextResult("", err), nil
		}

		profileName, err := s.getProfileName(browser, ctr, "visited pages from search engine query")
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
}
