package files

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	CoreDataOrigin = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
)

func getDb(path string) (*sql.DB, error) {
	return sql.Open("sqlite", fmt.Sprintf("file:%s?immutable=1", path))
}

func toDbDate(d time.Time) float64 {
	return float64(d.Unix() - CoreDataOrigin)
}

func fromDbDate(dbDate float64) time.Time {
	return time.Unix(int64(dbDate)+CoreDataOrigin, 0)
}
