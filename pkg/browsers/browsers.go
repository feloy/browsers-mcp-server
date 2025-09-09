package browsers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

var providers = map[string]api.Browser{}

// Register a new browser provider
func Register(provider api.Browser) {
	if provider == nil {
		panic("cannot register a nil browser provider")
	}
	providers[provider.Name()] = provider
}

// Clear the registered tools providers (Exposed for testing purposes)
func Clear() {
	providers = map[string]api.Browser{}
}

func GetBrowsers() []api.Browser {
	availableBrowsers := []api.Browser{}
	for _, provider := range providers {
		available, err := provider.IsAvailable()
		if err != nil {
			continue
		}
		if !available {
			continue
		}
		availableBrowsers = append(availableBrowsers, provider)
	}
	slices.SortFunc(availableBrowsers, func(a, b api.Browser) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return availableBrowsers
}

func GetBrowserByName(name string) (api.Browser, error) {
	found := false
	provider, found := providers[name]
	if !found {
		return nil, fmt.Errorf("browser %q not found", name)
	}
	available, err := provider.IsAvailable()
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, fmt.Errorf("browser %q is not available", name)
	}
	return provider, nil
}
