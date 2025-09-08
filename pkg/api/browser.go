package api

type Browser interface {
	Name() string
	IsAvailable() (bool, error)
	Profiles() ([]string, error)
}
