package files

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/feloy/mcp-server/pkg/system"
)

type Profile struct {
	ID         int
	Name       string
	IsRelative bool
	Path       string
	Default    bool
}

func ReadProfilesIni() ([]Profile, error) {
	path := filepath.Join(getUserDataDirecory(), "profiles.ini")
	data, err := system.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	var profiles []Profile

	sections := f.Sections()
	for _, section := range sections {
		name := section.Name()
		if strings.HasPrefix(name, "Profile") {
			profile, err := readProfile(section, strings.TrimPrefix(name, "Profile"))
			if err != nil {
				return nil, err
			}
			profiles = append(profiles, *profile)
		}
	}
	return profiles, nil
}

func readProfile(section *ini.Section, id string) (*Profile, error) {
	n, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("firefox: invalid profile ID: %w", err)
	}
	return &Profile{
		ID:         n,
		Name:       section.Key("Name").String(),
		IsRelative: section.Key("IsRelative").MustBool(),
		Path:       section.Key("Path").String(),
		Default:    section.Key("Default").MustBool(),
	}, nil
}
