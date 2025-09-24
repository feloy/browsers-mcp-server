package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

func (s *Server) initSourceReposVisits() []server.ServerTool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List the source repositories pages visited in the browser"),
	}

	browserProfiles := BrowsersProfiles{}
	browserProfiles.Populate(browsers.GetBrowsers())
	profilesEnum := browserProfiles.FlatList()
	log.Debug("source repos visits", "profilesEnum", profilesEnum)

	if len(profilesEnum) > 0 {
		options = append(options,
			mcp.WithString(
				"profile",
				mcp.Required(),
				mcp.Enum(profilesEnum...),
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
			Handler: s.listSourceReposVisits,
		},
	}

}

func (s *Server) listSourceReposVisits(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
