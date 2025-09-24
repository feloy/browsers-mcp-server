package files

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
)

func fromDbDate(dbDate int64) time.Time {
	return time.Unix(dbDate/1_000_000, 0)
}
func toDbDate(d time.Time) int64 {
	return d.Unix() * 1_000_000
}

func getDb(profile string, isRelative bool) (*sql.DB, error) {
	if isRelative {
		profile = filepath.Join(getUserDataDirecory(), profile)
	}
	path := filepath.Join(profile, "places.sqlite")
	return sql.Open("sqlite", fmt.Sprintf("file:%s?immutable=1", path))
}
