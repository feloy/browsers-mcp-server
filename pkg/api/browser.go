package api

type BookMark struct {
	Name   string
	URL    string
	Folder []string
}

type Browser interface {
	Name() string
	IsAvailable() (bool, error)
	Profiles() ([]string, error)
	Bookmarks(profile string) ([]BookMark, error)
}
