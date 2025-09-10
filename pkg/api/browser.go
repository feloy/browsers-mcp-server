package api

import "time"

type BookMark struct {
	Name            string    `yaml:"name"`
	URL             string    `yaml:"url"`
	Folder          []string  `yaml:"folder"`
	DateAdded       time.Time `yaml:"date_added,omitempty"`
	DateModified    time.Time `yaml:"date_modified,omitempty"`
	DateLastVisited time.Time `yaml:"date_last_visited,omitempty"`
}

type Browser interface {
	Name() string
	IsAvailable() (bool, error)
	Profiles() ([]string, error)
	Bookmarks(profile string) ([]BookMark, error)
}
