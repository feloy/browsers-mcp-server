package files

import (
	"os"
	"path/filepath"
	"time"

	"howett.net/plist"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

type Bookmark struct {
	Type          string        `plist:"WebBookmarkType"`
	Title         string        `plist:"Title"`
	Children      []Bookmark    `plist:"Children"`
	URIDictionary URIDictionary `plist:"URIDictionary"`
	URLString     string        `plist:"URLString"`
	DateAdded     time.Time     `plist:"dateAdded"`
}

type URIDictionary struct {
	Title string `plist:"title"`
}

func ListBookmarks() ([]api.BookMark, error) {
	path := filepath.Join(os.Getenv("HOME"), "Library", "Safari", "Bookmarks.plist")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var bookmarks Bookmark
	plist.NewDecoder(file).Decode(&bookmarks)
	return flatBookmarksRec(bookmarks, []string{}), nil
}

func flatBookmarksRec(bookmark Bookmark, folder []string) []api.BookMark {
	result := []api.BookMark{}
	if bookmark.Type == "WebBookmarkTypeList" {
		for _, child := range bookmark.Children {
			newFolder := folder
			if bookmark.Title != "" {
				newFolder = append(folder, bookmark.Title)
			}
			result = append(result, flatBookmarksRec(child, newFolder)...)
		}
	}
	if bookmark.Type == "WebBookmarkTypeLeaf" {
		result = append(result, api.BookMark{
			Name:      bookmark.URIDictionary.Title,
			URL:       bookmark.URLString,
			Folder:    folder,
			DateAdded: bookmark.DateAdded,
		})
	}
	return result
}
