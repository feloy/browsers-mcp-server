package api

import "time"

type BookMark struct {
	Name            string    `yaml:"name"`
	URL             string    `yaml:"url"`
	Folder          []string  `yaml:"folder"`
	DateAdded       time.Time `yaml:"date_added,omitempty"`
	DateModified    time.Time `yaml:"date_modified,omitempty"`
	DateLastVisited time.Time `yaml:"date_last_visited,omitempty"`
}

type SearchEngineQuery struct {
	Query        string    `yaml:"query"`
	Date         time.Time `yaml:"date"`
	SearchEngine string    `yaml:"search_engine"`
}

type SearchEngineOptions struct {
	StartTime time.Time
	EndTime   time.Time
	Limit     int
}

type VisitedPageFromSearchEngineQuery struct {
	URL          string    `yaml:"url"`
	Title        string    `yaml:"title"`
	Date         time.Time `yaml:"date"`
	SearchEngine string    `yaml:"search_engine"`
}

type ListVisitedPagesFromSearchEngineQueryOptions struct {
	Query     string
	StartTime time.Time
	EndTime   time.Time
}

type SourceRepoPageType string

const (
	SourceRepoPageTypeUnset            SourceRepoPageType = ""
	SourceRepoPageTypeProviderHome     SourceRepoPageType = "provider home"
	SourceRepoPageTypeOrganizationHome SourceRepoPageType = "organization home"
	SourceRepoPageTypeRepositoryHome   SourceRepoPageType = "repository home"
	SourceRepoPageTypeIssuesList       SourceRepoPageType = "issues list"
	SourceRepoPageTypePullRequestsList SourceRepoPageType = "pull requests list"
	SourceRepoPageTypeDiscussionsList  SourceRepoPageType = "discussions list"
	SourceRepoPageTypeIssue            SourceRepoPageType = "issue"
	SourceRepoPageTypePullRequest      SourceRepoPageType = "pull request"
	SourceRepoPageTypeDiscussion       SourceRepoPageType = "discussion"
	SourceRepoPageTypeOtherDetails     SourceRepoPageType = "other details"
)

type VisitedPageFromSourceRepos struct {
	Times        int                `yaml:"times"`
	Provider     string             `yaml:"provider"` // github, ...
	URL          string             `yaml:"url"`
	Organization string             `yaml:"organization"`
	Repository   string             `yaml:"repository"`
	Type         SourceRepoPageType `yaml:"type"`             // provider home, issue, pull request, ...
	Number       *string            `yaml:"number,omitempty"` // Depending on type: number of issue/PR/etc, not defined for home
}

type ListVisitedPagesFromSourceReposOptions struct {
	Type      SourceRepoPageType
	StartTime time.Time
	EndTime   time.Time
}

type Browser interface {
	Name() string
	IsAvailable() (bool, error)
	Profiles() ([]string, error)
	Bookmarks(profile string) ([]BookMark, error)
	SearchEngineQueries(profile string, options SearchEngineOptions) ([]SearchEngineQuery, error)
	ListVisitedPagesFromSearchEngineQuery(profile string, options ListVisitedPagesFromSearchEngineQueryOptions) ([]VisitedPageFromSearchEngineQuery, error)
	ListVisitedPagesFromSourceRepos(profile string, options ListVisitedPagesFromSourceReposOptions) ([]VisitedPageFromSourceRepos, error)
}
