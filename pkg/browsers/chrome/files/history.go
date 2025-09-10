package files

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

type queryResult struct {
	VisitTime int64
	URL       string
}

func SearchEngineQueries(profile string, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	filename := filepath.Join(getUserDataDirecory(), profile, "History")
	db, err := getDb(filename)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()))
	if options.StartTime != nil {
		startTime = toDbDate(*options.StartTime)
	}
	rows, err := db.Query(`SELECT 
	visits.visit_time,
	urls.url
FROM urls
INNER JOIN visits ON visits.url = urls.id
WHERE 
	urls.url like 'https://www.google.com/search%'
	AND visits.visit_time >= ?
	ORDER BY visits.visit_time ASC
LIMIT ?`, startTime, options.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchEngineQueries []api.SearchEngineQuery
	for rows.Next() {
		var queryResult queryResult
		err = rows.Scan(&queryResult.VisitTime, &queryResult.URL)
		if err != nil {
			return nil, err
		}
		urlParts, err := url.Parse(queryResult.URL)
		if err != nil {
			return nil, err
		}
		query := urlParts.Query().Get("q")

		searchEngineQueries = append(searchEngineQueries, api.SearchEngineQuery{
			Query:        query,
			Date:         fromDbDate(queryResult.VisitTime),
			SearchEngine: "Google",
		})
	}
	return searchEngineQueries, nil
}

func getDb(filename string) (*sql.DB, error) {
	return sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&nolock=1", filename))
}

func fromDbDate(dbDate int64) time.Time {
	return time.Unix(dbDate/1_000_000-11_644_473_600, 0)
}
func toDbDate(d time.Time) int64 {
	return (d.Unix() + 11_644_473_600) * 1_000_000
}
