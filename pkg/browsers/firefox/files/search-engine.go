package files

import (
	"net/url"

	"github.com/feloy/browsers-mcp-server/pkg/api"
)

func SearchEngineQueries(profile string, isRelative bool, options api.SearchEngineOptions) ([]api.SearchEngineQuery, error) {
	type queryResult struct {
		VisitDate int64
		URL       string
	}

	db, err := getDb(profile, isRelative)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(options.StartTime)
	endTime := toDbDate(options.EndTime)
	rows, err := db.Query(`SELECT
  visit_date,
	url 
FROM moz_historyvisits hv
INNER JOIN moz_places p ON p.id = hv.place_id 
WHERE url LIKE 'https://www.google.com/search%'
AND hv.visit_date >= ?
AND hv.visit_date < ?
ORDER BY hv.visit_date ASC
LIMIT ?`, startTime, endTime, options.Limit)
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

func ListVisitedPagesFromSearchEngineQuery(profile string, isRelative bool, options api.ListVisitedPagesFromSearchEngineQueryOptions) ([]api.VisitedPageFromSearchEngineQuery, error) {
	type queryResult struct {
		VisitTime int64
		URL       string
		Title     string
	}

	db, err := getDb(profile, isRelative)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(options.StartTime)
	endTime := toDbDate(options.EndTime)
	rows, err := db.Query(`SELECT
  visited.visit_date,
	visited_place.url,
	visited_place.title
FROM moz_historyvisits hv
INNER JOIN moz_places p ON p.id = hv.place_id 
INNER JOIN moz_historyvisits visited ON visited.from_visit = hv.id
INNER JOIN moz_places visited_place ON visited_place.id = visited.place_id
WHERE p.url LIKE 'https://www.google.com/search%'
AND (? = '' OR p.url like ? OR p.url like ?)
AND hv.visit_date >= ?
AND hv.visit_date < ?
ORDER BY hv.visit_date ASC`, options.Query, "%q="+url.QueryEscape(options.Query)+"&%", "%q="+url.QueryEscape(options.Query), startTime, endTime)
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
