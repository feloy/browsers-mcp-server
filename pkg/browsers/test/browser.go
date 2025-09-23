package test

import "github.com/feloy/browsers-mcp-server/pkg/api"

var _ api.Browser = &Browser{}

type Browser struct {
	name                                   string
	available                              bool
	availableError                         error
	profiles                               []string
	profilesError                          error
	bookmarks                              []api.BookMark
	bookmarksError                         error
	searchEngineQueries                    []api.SearchEngineQuery
	searchEngineQueriesError               error
	visitedPagesFromSearchEngineQuery      []api.VisitedPageFromSearchEngineQuery
	visitedPagesFromSearchEngineQueryError error
	visitedPagesFromSourceRepos            []api.VisitedPageFromSourceRepos
	visitedPagesFromSourceReposError       error
}

type NewBrowserOptions struct {
	Name                                   string
	Available                              bool
	AvailableError                         error
	Profiles                               []string
	ProfilesError                          error
	Bookmarks                              []api.BookMark
	BookmarksError                         error
	SearchEngineQueries                    []api.SearchEngineQuery
	SearchEngineQueriesError               error
	VisitedPagesFromSearchEngineQuery      []api.VisitedPageFromSearchEngineQuery
	VisitedPagesFromSearchEngineQueryError error
	VisitedPagesFromSourceRepos            []api.VisitedPageFromSourceRepos
	VisitedPagesFromSourceReposError       error
}

func NewBrowser(options NewBrowserOptions) *Browser {
	return &Browser{
		name:                                   options.Name,
		available:                              options.Available,
		availableError:                         options.AvailableError,
		profiles:                               options.Profiles,
		profilesError:                          options.ProfilesError,
		bookmarks:                              options.Bookmarks,
		bookmarksError:                         options.BookmarksError,
		searchEngineQueries:                    options.SearchEngineQueries,
		searchEngineQueriesError:               options.SearchEngineQueriesError,
		visitedPagesFromSearchEngineQuery:      options.VisitedPagesFromSearchEngineQuery,
		visitedPagesFromSearchEngineQueryError: options.VisitedPagesFromSearchEngineQueryError,
		visitedPagesFromSourceRepos:            options.VisitedPagesFromSourceRepos,
		visitedPagesFromSourceReposError:       options.VisitedPagesFromSourceReposError,
	}
}

func (o *Browser) Name() string {
	return o.name
}

func (o *Browser) IsAvailable() (bool, error) {
	return o.available, o.availableError
}

func (o *Browser) Profiles() ([]string, error) {
	return o.profiles, o.profilesError
}

func (o *Browser) Bookmarks(profile string) ([]api.BookMark, error) {
	return o.bookmarks, o.bookmarksError
}

func (o *Browser) SearchEngineQueries(profile string, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	return o.searchEngineQueries, o.searchEngineQueriesError
}

func (o *Browser) ListVisitedPagesFromSearchEngineQuery(profile string, options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	return o.visitedPagesFromSearchEngineQuery, o.visitedPagesFromSearchEngineQueryError
}

func (o *Browser) ListVisitedPagesFromSourceRepos(profile string, options api.ListVisitedPagesFromSourceReposOptions) ([]api.VisitedPageFromSourceRepos, error) {
	return o.visitedPagesFromSourceRepos, o.visitedPagesFromSourceReposError
}
