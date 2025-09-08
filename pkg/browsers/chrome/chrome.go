package chrome

import (
	"github.com/feloy/mcp-server/pkg/api"
	"github.com/feloy/mcp-server/pkg/browsers"
	"github.com/feloy/mcp-server/pkg/browsers/chrome/files"
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

func init() {
	browsers.Register(instance)
}
