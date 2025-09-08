package test

import "github.com/feloy/mcp-server/pkg/api"

var _ api.Browser = &Browser{}

type Browser struct {
	name           string
	available      bool
	availableError error
	profiles       []string
	profilesError  error
}

func NewBrowser(name string, available bool, availableError error, profiles []string, profilesError error) *Browser {
	return &Browser{
		name:           name,
		available:      available,
		availableError: availableError,
		profiles:       profiles,
		profilesError:  profilesError,
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
