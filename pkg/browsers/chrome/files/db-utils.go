package files

import (
	"database/sql"
	"fmt"
	"time"
)

func getDb(filename string) (*sql.DB, error) {
	return sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&nolock=1", filename))
}

func fromDbDate(dbDate int64) time.Time {
	return time.Unix(dbDate/1_000_000-11_644_473_600, 0)
}
func toDbDate(d time.Time) int64 {
	return (d.Unix() + 11_644_473_600) * 1_000_000
}
