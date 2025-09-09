package files

import (
	"database/sql"
	"path/filepath"

	"github.com/feloy/browsers-mcp-server/pkg/api"
	_ "modernc.org/sqlite"
)

func ListBookmarks(profile string, isRelative bool) ([]api.BookMark, error) {
	result := []api.BookMark{}
	db, err := getDb(profile, isRelative)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = listBookmarksRec(db, 0, []string{}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func listBookmarksRec(db *sql.DB, parent int, folder []string, result *[]api.BookMark) error {
	subdirs, titles, err := getSubdirs(db, parent)
	if err != nil {
		return err
	}

	leafs, leafTitles, leafUrls, err := getLeafs(db, parent)
	if err != nil {
		return err
	}
	if len(leafs) > 0 {
		for i := range leafs {
			*result = append(*result, api.BookMark{
				Name:   leafTitles[i],
				URL:    leafUrls[i],
				Folder: folder,
			})
		}
	}

	for i, subdir := range subdirs {
		newFolder := folder
		if titles[i] != "" {
			newFolder = append(folder, titles[i])
		}
		listBookmarksRec(db, subdir, newFolder, result)
	}
	return nil
}

func getDb(profile string, isRelative bool) (*sql.DB, error) {
	if isRelative {
		profile = filepath.Join(getUserDataDirecory(), profile)
	}
	path := filepath.Join(profile, "places.sqlite")
	return sql.Open("sqlite", path)
}

func getSubdirs(db *sql.DB, parent int) ([]int, []string, error) {
	rows, err := db.Query("SELECT id, title FROM moz_bookmarks WHERE parent = ? AND type = 2", parent)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var ids []int
	var titles []string
	for rows.Next() {
		var id int
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		titles = append(titles, title)
	}
	return ids, titles, nil
}

func getLeafs(db *sql.DB, parent int) ([]int, []string, []string, error) {
	rows, err := db.Query(" SELECT b.id, b.title, p.url FROM moz_bookmarks AS b INNER JOIN moz_places AS p ON b.fk = p.id WHERE b.parent = ? AND b.type = 1;", parent)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	var ids []int
	var titles []string
	var urls []string
	for rows.Next() {
		var id int
		var title string
		var url string
		err = rows.Scan(&id, &title, &url)
		if err != nil {
			return nil, nil, nil, err
		}
		ids = append(ids, id)
		titles = append(titles, title)
		urls = append(urls, url)
	}
	return ids, titles, urls, nil
}
