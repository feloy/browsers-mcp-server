package files

import (
	"os"
	"path/filepath"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

func ListBookmarks() ([]api.BookMark, error) {
	path := filepath.Join(os.Getenv("HOME"), "Library", "Safari", "Bookmarks.plist")
	_ = path
	return []api.BookMark{}, nil
}
