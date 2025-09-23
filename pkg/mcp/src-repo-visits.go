package mcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initSourceReposVisits() []server.ServerTool {
	browsers := browsers.GetBrowsers()
	if len(browsers) == 1 {
		options := []mcp.ToolOption{
			mcp.WithDescription("List the source repositories pages visited in the browser"),
		}
		profiles, err := browsers[0].Profiles()
		if err != nil {
			return nil
		}
		if len(profiles) > 1 {
			options = append(options,
				mcp.WithString(
					"profile",
					mcp.Required(),
					mcp.Enum(profiles...),
					mcp.Description("The browser's profile to list the visits for"),
				))
		}
		options = append(
			options,
			mcp.WithString(
				"day",
				mcp.Description("List the visits done on this day (YYYY-MM-DD), default is today"),
			),
			mcp.WithString(
				"type",
				mcp.Description("The type of source repository page to list. If not set, returns all pages of all types."),
				mcp.Enum(
					string(api.SourceRepoPageTypeProviderHome),
					string(api.SourceRepoPageTypeOrganizationHome),
					string(api.SourceRepoPageTypeRepositoryHome),
					string(api.SourceRepoPageTypeIssuesList),
					string(api.SourceRepoPageTypePullRequestsList),
					string(api.SourceRepoPageTypeDiscussionsList),
					string(api.SourceRepoPageTypeIssue),
					string(api.SourceRepoPageTypePullRequest),
					string(api.SourceRepoPageTypeDiscussion),
					string(api.SourceRepoPageTypeOtherDetails),
				),
			),
			mcp.WithString(
				"name",
				mcp.Description("The name of the source repository page to list"),
			),
		)
		return []server.ServerTool{
			{
				Tool:    mcp.NewTool("list_source_repos_visits", options...),
				Handler: s.listSourceReposVisitsForBrowser(nil),
			},
		}
	} else {
		var tools []server.ServerTool
		for _, browser := range browsers {
			options := []mcp.ToolOption{
				mcp.WithDescription(fmt.Sprintf("List the source repositories pages visited in browser %s", browser.Name())),
			}
			profiles, err := browser.Profiles()
			if err != nil {
				return nil
			}
			if len(profiles) > 1 {
				options = append(options,
					mcp.WithString(
						"profile",
						mcp.Required(),
						mcp.Enum(profiles...),
						mcp.Description("The browser's profile to list the bookmarks for"),
					),
				)
			}
			options = append(
				options,
				mcp.WithString(
					"day",
					mcp.Description("List the visits done on this day (YYYY-MM-DD), default is today"),
				),
				mcp.WithString(
					"type",
					mcp.Description("The type of source repository page to list. If not set, returns all pages of all types."),
					mcp.Enum(
						string(api.SourceRepoPageTypeProviderHome),
						string(api.SourceRepoPageTypeOrganizationHome),
						string(api.SourceRepoPageTypeRepositoryHome),
						string(api.SourceRepoPageTypeIssuesList),
						string(api.SourceRepoPageTypePullRequestsList),
						string(api.SourceRepoPageTypeDiscussionsList),
						string(api.SourceRepoPageTypeIssue),
						string(api.SourceRepoPageTypePullRequest),
						string(api.SourceRepoPageTypeDiscussion),
						string(api.SourceRepoPageTypeOtherDetails),
					),
				))
			browserName := browser.Name()
			tools = append(tools, server.ServerTool{
				Tool:    mcp.NewTool(fmt.Sprintf("list_source_repos_visits_%s", browser.Name()), options...),
				Handler: s.listSourceReposVisitsForBrowser(&browserName),
			})
		}
		return tools
	}
}

func (s *Server) listSourceReposVisitsForBrowser(browserName *string) func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		profileName, err := s.getProfileName(browser, ctr, "list of source repository visits")
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

		var pageType api.SourceRepoPageType
		if pageTypeStr, ok := ctr.GetArguments()["type"].(string); ok {
			pageType = api.SourceRepoPageType(pageTypeStr)
		}

		visits, err := browser.ListVisitedPagesFromSourceRepos(profileName, api.ListVisitedPagesFromSourceReposOptions{
			Type:      pageType,
			StartTime: startTime,
			EndTime:   endTime,
		})
		if err != nil {
			return NewTextResult("", err), nil
		}

		yamlVisits, err := yaml.Marshal(visits)
		if err != nil {
			return NewTextResult("", err), nil
		}
		return NewTextResult(fmt.Sprintf("The following visits to source repositories (YAML format) were found:\n%s", string(yamlVisits)), nil), nil
	}
}
