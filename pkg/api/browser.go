package api

import "time"

type BookMark struct {
	Name            string    `yaml:"name" json:"bookmark_name"`
	URL             string    `yaml:"url" json:"bookmark_url"`
	Folder          []string  `yaml:"folder" json:"bookmark_folder"`
	DateAdded       time.Time `yaml:"date_added,omitempty" json:"bookmark_date_added"`
	DateModified    time.Time `yaml:"date_modified,omitempty" json:"bookmark_date_modified"`
	DateLastVisited time.Time `yaml:"date_last_visited,omitempty" json:"bookmark_date_last_visited"`
}

type SearchEngineQuery struct {
	Query        string    `yaml:"query" json:"query"`
	Date         time.Time `yaml:"date" json:"date"`
	SearchEngine string    `yaml:"search_engine" json:"search_engine"`
}

type SearchEngineOptions struct {
	StartTime *time.Time
	Limit     int
}

type VisitedPageFromSearchEngineQuery struct {
	URL          string    `yaml:"url"`
	Title        string    `yaml:"title"`
	Date         time.Time `yaml:"date"`
	SearchEngine string    `yaml:"search_engine"`
}

type ListVisitedPagesFromSearchEngineQueryOptions struct {
	Query     string
	StartTime *time.Time
}

type Browser interface {
	Name() string
	IsAvailable() (bool, error)
	Profiles() ([]string, error)
	Bookmarks(profile string) ([]BookMark, error)
	SearchEngineQueries(profile string, options SearchEngineOptions) ([]SearchEngineQuery, error)
	ListVisitedPagesFromSearchEngineQuery(profile string, options ListVisitedPagesFromSearchEngineQueryOptions) ([]VisitedPageFromSearchEngineQuery, error)
}
