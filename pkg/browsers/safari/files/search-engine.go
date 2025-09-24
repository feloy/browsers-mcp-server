package files

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

func SearchEngineQueries(options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	type queryResult struct {
		VisitTime float64
		URL       string
	}

	path := filepath.Join(os.Getenv("HOME"), "Library", "Safari", "History.db")
	db, err := getDb(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(options.StartTime)
	endTime := toDbDate(options.EndTime)
	rows, err := db.Query(`SELECT DISTINCT
	round(history_visits.visit_time),
	history_items.url
FROM history_visits
INNER JOIN history_items ON history_items.id = history_visits.history_item
WHERE 
	history_items.url like 'https://www.google.com/search%'
	AND visit_time >= ?
	AND visit_time < ?
	ORDER BY visit_time ASC
LIMIT ?`, startTime, endTime, options.Limit)
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

func ListVisitedPagesFromSearchEngineQuery(options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	return []api.VisitedPageFromSearchEngineQuery{}, nil
}
