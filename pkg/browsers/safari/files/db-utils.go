package files

import (
	"database/sql"
	"fmt"
)

func getDb(path string) (*sql.DB, error) {
	return sql.Open("sqlite", fmt.Sprintf("file:%s?immutable=1", path))
}
