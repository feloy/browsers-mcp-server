package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/feloy/mcp-server/pkg/system"
	"github.com/spf13/afero"
)

func TestReadProfilesIni(t *testing.T) {
	system.FileSystem = afero.NewMemMapFs()
	basePath := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Firefox")
	system.WriteFile(filepath.Join(basePath, "profiles.ini"), []byte(`[Profile0]
Name=my-profile
IsRelative=1
Path=path/to/profile
Default=1
`), 0644)
	profiles, err := ReadProfilesIni()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(profiles) == 0 {
		t.Errorf("Expected at least one profile, got %d", len(profiles))
	}
	if profiles[0].Name != "my-profile" {
		t.Errorf("Expected my-profile, got %s", profiles[0].Name)
	}
	if profiles[0].IsRelative != true {
		t.Errorf("Expected IsRelative to be true, got %v", profiles[0].IsRelative)
	}
	if profiles[0].Path != "path/to/profile" {
		t.Errorf("Expected path/to/profile, got %s", profiles[0].Path)
	}
	if profiles[0].Default != true {
		t.Errorf("Expected Default to be true, got %v", profiles[0].Default)
	}
}
