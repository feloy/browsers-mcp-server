package files

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/system"
	globaltest "github.com/feloy/browsers-mcp-server/pkg/test"
	"github.com/spf13/afero"
)

func TestListBookmarks(t *testing.T) {
	system.FileSystem = afero.NewMemMapFs()
	system.Os = "darwin"
	basePath := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Google", "Chrome", "Profile1")
	system.WriteFile(filepath.Join(basePath, "Bookmarks"), []byte(`{
  "roots": {
		"synced": {
       "children": [{
				"name": "RedHat",
				"type": "url",
				"url": "https://www.redhat.com",
				"date_added": "13390300334000000",
				"date_modified": "13390300335000000",
				"date_last_used": "13390300336000000"
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
	if bookmarks[0].DateAdded != globaltest.Must(time.Parse(time.RFC3339, "2025-04-28T07:52:14Z")) {
		t.Errorf("Expected 2021-01-01T00:00:00Z, got %s", bookmarks[0].DateAdded)
	}
	if bookmarks[0].DateModified != globaltest.Must(time.Parse(time.RFC3339, "2025-04-28T07:52:15Z")) {
		t.Errorf("Expected 2021-01-01T00:00:00Z, got %s", bookmarks[0].DateModified)
	}
	if bookmarks[0].DateLastVisited != globaltest.Must(time.Parse(time.RFC3339, "2025-04-28T07:52:16Z")) {
		t.Errorf("Expected 2021-01-01T00:00:00Z, got %s", bookmarks[0].DateLastVisited)
	}
}
