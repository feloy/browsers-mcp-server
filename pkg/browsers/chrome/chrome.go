package chrome

import (
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/chrome/files"
)

var instance api.Browser = &Chrome{}

type Chrome struct{}

func (o *Chrome) Name() string {
	return "chrome"
}

func (o *Chrome) IsAvailable() (bool, error) {
	_, err := files.ReadLocalState()
	return err == nil, err
}

func (o *Chrome) Profiles() ([]string, error) {
	localState, err := files.ReadLocalState()
	if err != nil {
		return nil, err
	}
	return localState.Profile.ProfilesOrder, nil
}

func (o *Chrome) Bookmarks(profileName string) ([]api.BookMark, error) {
	profiles, err := o.Profiles()
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile == profileName {
			return files.ListBookmarks(profile)
		}
	}
	return nil, fmt.Errorf("profile %s not found", profileName)
}

func (o *Chrome) SearchEngineQueries(profileName string, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	profiles, err := o.Profiles()
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile == profileName {
			return files.SearchEngineQueries(profile, options)
		}
	}
	return nil, fmt.Errorf("profile %s not found", profileName)
}

func (o *Chrome) ListVisitedPagesFromSearchEngineQuery(profileName string, options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	profiles, err := o.Profiles()
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile == profileName {
			return files.ListVisitedPagesFromSearchEngineQuery(profile, options)
		}
	}
	return nil, fmt.Errorf("profile %s not found", profileName)
}

func init() {
	browsers.Register(instance)
}
