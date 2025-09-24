package mcp

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/api"
)

// key: browser name, value: profiles names
type BrowsersProfiles map[string][]string

func (b *BrowsersProfiles) Populate(browsers []api.Browser) {
	for _, browser := range browsers {
		profiles, err := browser.Profiles()
		if err != nil {
			log.Error("failed to get profiles for browser", "browser", browser.Name(), "error", err)
			continue
		}
		(*b)[browser.Name()] = profiles
	}
}

func (b *BrowsersProfiles) FlatList() []string {
	if len(*b) == 0 {
		// no browsers found
		return []string{}
	}
	if len(*b) == 1 && len(slices.Collect(maps.Values(*b))[0]) == 1 {
		// only one browser found, with a single profile
		return []string{}
	}

	multipleBrowsers := len(*b) > 1

	result := []string{}
	for browserName, profiles := range *b {
		for _, profile := range profiles {
			if len(profiles) > 1 {
				if multipleBrowsers {
					result = append(result, fmt.Sprintf("%s on %s", profile, browserName))
				} else {
					result = append(result, profile)
				}
			} else {
				result = append(result, browserName)
			}
		}
	}
	return result
}

func GetBrowserAndProfileFromValue(value string, browsers []api.Browser) (string, string, error) {
	browserProfiles := map[string][]string{}
	for _, browser := range browsers {
		profiles, err := browser.Profiles()
		if err != nil {
			continue
		}
		browserProfiles[browser.Name()] = profiles
	}

	parts := strings.Split(value, " on ")

	if len(parts) == 2 {
		var profiles []string
		var ok bool
		if profiles, ok = browserProfiles[parts[1]]; !ok {
			return "", "", fmt.Errorf("browser %q not found", parts[1])
		}
		if !slices.Contains(profiles, parts[0]) {
			return "", "", fmt.Errorf("profile %q not found", parts[0])
		}
		return parts[1], parts[0], nil
	}

	// No value: should be one browser with one profile
	if len(parts) == 1 && len(parts[0]) == 0 {
		if len(browserProfiles) == 1 {
			browserName := slices.Collect(maps.Keys(browserProfiles))[0]
			profiles := browserProfiles[browserName]
			if len(profiles) == 1 {
				return browserName, profiles[0], nil
			}
			return "", "", fmt.Errorf("multiple profiles found for browser %q, value cannot be empty", browserName)
		}
		return "", "", fmt.Errorf("multiple browsers found, value cannot be empty")
	}

	// Single value

	// - first check if this is the name of a browser having a single profile
	if profiles, ok := browserProfiles[value]; ok && len(profiles) == 1 {
		return value, profiles[0], nil
	}

	// - then check if this the name of a profile of the single browser
	if len(browserProfiles) == 1 {
		browserName := slices.Collect(maps.Keys(browserProfiles))[0]
		profiles := browserProfiles[browserName]
		if slices.Contains(profiles, value) {
			return browserName, value, nil
		}
	}

	return "", "", errors.New("incorrect profile or browser name")
}
