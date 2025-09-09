package files

import (
	"encoding/json"
	"path/filepath"
	"slices"

	"github.com/andrewarchi/browser/jsonutil/timefmt"
	"github.com/andrewarchi/browser/jsonutil/uuid"
	"github.com/feloy/browsers-mcp-server/pkg/api"
	"github.com/feloy/browsers-mcp-server/pkg/browsers/chrome/files/fields"
	"github.com/feloy/browsers-mcp-server/pkg/system"
)

// Bookmarks contains Chrome bookmark information.
type Bookmarks struct {
	Roots BookmarkRoots `json:"roots"`
}

// BookmarkRoots contains the root level bookmarks folders.
type BookmarkRoots struct {
	BookmarkBar BookmarkEntry `json:"bookmark_bar"` // "Bookmarks" folder
	Other       BookmarkEntry `json:"other"`        // "Other Bookmarks" folder
	Synced      BookmarkEntry `json:"synced"`       // "Mobile Bookmarks" folder
}

type BookmarkEntry struct {
	Children     []BookmarkEntry     `json:"children"` // for folder type only
	DateAdded    fields.QuotedChrome `json:"date_added"`
	DateModified fields.QuotedChrome `json:"date_modified,omitempty"` // for folder type only
	GUID         *uuid.UUID          `json:"guid"`                    // "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	ID           string              `json:"id"`                      // e.g. "567"
	Name         string              `json:"name"`
	Type         string              `json:"type"` // "folder" or "url"
	MetaInfo     *BookmarkMetaInfo   `json:"meta_info,omitempty"`
	URL          string              `json:"url,omitempty"` // for url type only
}

// BookmarkMetaInfo contains additional bookmark metadata.
type BookmarkMetaInfo struct {
	LastVisitedDesktop timefmt.QuotedChrome `json:"last_visited_desktop"`
}

// ParseBookmarks returns the bookmarks in a Chrome profile.
func ListBookmarks(profile string) ([]api.BookMark, error) {
	filename := filepath.Join(getUserDataDirecory(), profile, "Bookmarks")
	data, err := system.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var treeBookmarks Bookmarks
	err = json.Unmarshal(data, &treeBookmarks)
	if err != nil {
		return nil, err
	}

	flatten := slices.Concat(
		flatBookmarksRec(treeBookmarks.Roots.BookmarkBar, []string{"toolbar"}),
		flatBookmarksRec(treeBookmarks.Roots.Other, []string{"other"}),
		flatBookmarksRec(treeBookmarks.Roots.Synced, []string{"synced"}),
	)
	return flatten, nil
}

func flatBookmarksRec(entry BookmarkEntry, folder []string) []api.BookMark {
	result := []api.BookMark{}
	if entry.Type == "folder" {
		for _, child := range entry.Children {
			result = append(result, flatBookmarksRec(child, append(folder, entry.Name))...)
		}
	}
	if entry.Type == "url" {
		result = append(result, api.BookMark{
			Name:   entry.Name,
			URL:    entry.URL,
			Folder: folder,
		})
	}
	return result
}
