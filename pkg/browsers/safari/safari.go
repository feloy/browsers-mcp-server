package safari

import (
	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/safari/files"
	"github.com/feloy/browsers-mcp-server/pkg/system"
)

var instance api.Browser = &Safari{}

type Safari struct{}

func (o *Safari) Name() string {
	return "safari"
}

func (o *Safari) IsAvailable() (bool, error) {
	return system.Os == "darwin", nil
}

func (o *Safari) Profiles() ([]string, error) {
	// TODO support multiple profiles
	return []string{"DefaultProfile"}, nil
}

func (o *Safari) Bookmarks(profileName string) ([]api.BookMark, error) {
	return files.ListBookmarks()
}

func (o *Safari) SearchEngineQueries(profileName string, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	return files.SearchEngineQueries(options)
}

func (o *Safari) ListVisitedPagesFromSearchEngineQuery(profileName string, options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	return files.ListVisitedPagesFromSearchEngineQuery(options)
}

func (o *Safari) ListVisitedPagesFromSourceRepos(profileName string, options api.ListVisitedPagesFromSourceReposOptions) ([]api.VisitedPageFromSourceRepos, error) {
	return files.ListVisitedPagesFromSourceRepos(options)
}

func init() {
	browsers.Register(instance)
}
