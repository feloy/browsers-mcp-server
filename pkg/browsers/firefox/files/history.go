package files

import (
	"net/url"
	"time"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

type queryResult struct {
	VisitDate int64
	URL       string
}

func SearchEngineQueries(profile string, isRelative bool, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	db, err := getDb(profile, isRelative)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()))
	if options.StartTime != nil {
		startTime = toDbDate(*options.StartTime)
	}
	rows, err := db.Query(`SELECT
  visit_date,
	url 
FROM moz_historyvisits hv
INNER JOIN moz_places p ON p.id = hv.place_id 
WHERE url LIKE 'https://www.google.com/search%'
AND hv.visit_date >= ?
ORDER BY hv.visit_date ASC
LIMIT ?`, startTime, options.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searchEngineQueries []api.SearchEngineQuery
	for rows.Next() {
		var queryResult queryResult
		err = rows.Scan(&queryResult.VisitDate, &queryResult.URL)
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
			Date:         fromDbDate(queryResult.VisitDate),
			SearchEngine: "Google",
		})
	}
	return searchEngineQueries, nil
}

func fromDbDate(dbDate int64) time.Time {
	return time.Unix(dbDate/1_000_000, 0)
}
func toDbDate(d time.Time) int64 {
	return d.Unix() * 1_000_000
}
