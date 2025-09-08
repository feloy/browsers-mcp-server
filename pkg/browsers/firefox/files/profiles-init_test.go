package files

import (
	"fmt"
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
	fmt.Printf("profiles: %v", profiles)
}
