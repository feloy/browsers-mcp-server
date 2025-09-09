package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/feloy/mcp-server/pkg/system"
	"github.com/spf13/afero"
)

func TestListBookmarks(t *testing.T) {
	system.FileSystem = afero.NewMemMapFs()
	basePath := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Google", "Chrome", "Profile1")
	system.WriteFile(filepath.Join(basePath, "Bookmarks"), []byte(`{
  "roots": {
		"synced": {
       "children": [{
				"name": "RedHat",
				"type": "url",
				"url": "https://www.redhat.com"
			}],
       "name": "Mobile Bookmarks",
       "type": "folder"
    }
  }
}
`), 0644)

	bookmarks, err := ListBookmarks("Profile1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(bookmarks))
	}

	if bookmarks[0].Name != "RedHat" {
		t.Errorf("Expected RedHat, got %s", bookmarks[0].Name)
	}
}
