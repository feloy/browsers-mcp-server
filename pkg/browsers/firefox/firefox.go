package chrome

import (
	"fmt"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/firefox/files"
)

var instance api.Browser = &Firefox{}

type Firefox struct{}

func (o *Firefox) Name() string {
	return "firefox"
}

func (o *Firefox) IsAvailable() (bool, error) {
	_, err := files.ReadProfilesIni()
	return err == nil, err
}

func (o *Firefox) Profiles() ([]string, error) {
	profiles, err := files.ReadProfilesIni()
	if err != nil {
		return nil, err
	}
	profileNames := []string{}
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name)
	}
	return profileNames, nil
}

func (o *Firefox) Bookmarks(profileName string) ([]api.BookMark, error) {
	profiles, err := files.ReadProfilesIni()
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile.Name == profileName {
			return files.ListBookmarks(profile.Path, profile.IsRelative)
		}
	}
	return nil, fmt.Errorf("profile %s not found", profileName)
}

func init() {
	browsers.Register(instance)
}
