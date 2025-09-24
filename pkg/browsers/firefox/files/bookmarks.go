package files

import (
	"database/sql"
	"time"

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
	subdirs, err := getSubdirs(db, parent)
	if err != nil {
		return err
	}

	leafs, err := getLeafs(db, parent)
	if err != nil {
		return err
	}
	if len(leafs) > 0 {
		for i := range leafs {
			bm := api.BookMark{
				Name:   leafs[i].title,
				URL:    leafs[i].url,
				Folder: folder,
			}
			if leafs[i].dateAdded != nil {
				bm.DateAdded = time.Unix(int64(*leafs[i].dateAdded/1_000_000), 0)
			}
			if leafs[i].dateModified != nil {
				bm.DateModified = time.Unix(int64(*leafs[i].dateModified/1_000_000), 0)
			}
			if leafs[i].dateLastVisited != nil {
				bm.DateLastVisited = time.Unix(int64(*leafs[i].dateLastVisited/1_000_000), 0)
			}
			*result = append(*result, bm)
		}
	}

	for i, subdir := range subdirs {
		newFolder := folder
		if subdirs[i].title != "" {
			newFolder = append(folder, subdirs[i].title)
		}
		listBookmarksRec(db, subdir.id, newFolder, result)
	}
	return nil
}

type subdir struct {
	id    int
	title string
}

func getSubdirs(db *sql.DB, parent int) ([]subdir, error) {
	rows, err := db.Query("SELECT id, title FROM moz_bookmarks WHERE parent = ? AND type = 2", parent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subdirs []subdir
	for rows.Next() {
		var subdir subdir
		err = rows.Scan(&subdir.id, &subdir.title)
		if err != nil {
			return nil, err
		}
		subdirs = append(subdirs, subdir)
	}
	return subdirs, nil
}

type leaf struct {
	id              int
	title           string
	url             string
	dateAdded       *int
	dateModified    *int
	dateLastVisited *int
}

func getLeafs(db *sql.DB, parent int) ([]leaf, error) {
	rows, err := db.Query("SELECT b.id, b.title, p.url, b.dateAdded, b.lastModified, p.last_visit_date FROM moz_bookmarks AS b INNER JOIN moz_places AS p ON b.fk = p.id WHERE b.parent = ? AND b.type = 1;", parent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leafs []leaf
	for rows.Next() {
		var leaf leaf
		err = rows.Scan(&leaf.id, &leaf.title, &leaf.url, &leaf.dateAdded, &leaf.dateModified, &leaf.dateLastVisited)
		if err != nil {
			return nil, err
		}
		leafs = append(leafs, leaf)
	}
	return leafs, nil
}
