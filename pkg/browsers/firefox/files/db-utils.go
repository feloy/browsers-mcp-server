package files

import (
	"time"
)

func fromDbDate(dbDate int64) time.Time {
	return time.Unix(dbDate/1_000_000, 0)
}
func toDbDate(d time.Time) int64 {
	return d.Unix() * 1_000_000
}
