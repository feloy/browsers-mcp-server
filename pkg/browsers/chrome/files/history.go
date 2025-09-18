package files

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/api"
)

func SearchEngineQueries(profile string, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	log.Debug("searching engine queries", "profile", profile, "options", options)

	type queryResult struct {
		VisitTime int64
		URL       string
	}

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

func ListVisitedPagesFromSearchEngineQuery(profile string, options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	type queryResult struct {
		VisitTime int64
		URL       string
		Title     string
	}
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
visited.visit_time,
visited_url.url,
visited_url.title
FROM urls
INNER JOIN visits ON visits.url = urls.id
INNER JOIN visits visited on visited.from_visit = visits.id
INNER JOIN urls visited_url on visited_url.id = visited.url
WHERE 
  urls.url like 'https://www.google.com/search%'
	AND (? = '' OR urls.url like ? OR urls.url like ?)
  AND visits.visit_time >= ?
ORDER BY visits.visit_time ASC`, options.Query, "%q="+url.QueryEscape(options.Query)+"&%", "%q="+url.QueryEscape(options.Query), startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visitedPages []api.VisitedPageFromSearchEngineQuery
	for rows.Next() {
		var queryResult queryResult
		err = rows.Scan(&queryResult.VisitTime, &queryResult.URL, &queryResult.Title)
		if err != nil {
			return nil, err
		}

		visitedPages = append(visitedPages, api.VisitedPageFromSearchEngineQuery{
			URL:          queryResult.URL,
			Title:        queryResult.Title,
			Date:         fromDbDate(queryResult.VisitTime),
			SearchEngine: "Google",
		})
	}
	return visitedPages, nil
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
