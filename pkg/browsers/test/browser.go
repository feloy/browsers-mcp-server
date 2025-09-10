package test

import "github.com/feloy/browsers-mcp-server/pkg/api"

var _ api.Browser = &Browser{}

type Browser struct {
	name                     string
	available                bool
	availableError           error
	profiles                 []string
	profilesError            error
	bookmarks                []api.BookMark
	bookmarksError           error
	searchEngineQueries      []api.SearchEngineQuery
	searchEngineQueriesError error
}

type NewBrowserOptions struct {
	Name                     string
	Available                bool
	AvailableError           error
	Profiles                 []string
	ProfilesError            error
	Bookmarks                []api.BookMark
	BookmarksError           error
	SearchEngineQueries      []api.SearchEngineQuery
	SearchEngineQueriesError error
}

func NewBrowser(options NewBrowserOptions) *Browser {
	return &Browser{
		name:                     options.Name,
		available:                options.Available,
		availableError:           options.AvailableError,
		profiles:                 options.Profiles,
		profilesError:            options.ProfilesError,
		bookmarks:                options.Bookmarks,
		bookmarksError:           options.BookmarksError,
		searchEngineQueries:      options.SearchEngineQueries,
		searchEngineQueriesError: options.SearchEngineQueriesError,
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
