package chrome

import (
	"github.com/feloy/mcp-server/pkg/api"
	"github.com/feloy/mcp-server/pkg/browsers"
	"github.com/feloy/mcp-server/pkg/browsers/firefox/files"
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

func init() {
	browsers.Register(instance)
}
