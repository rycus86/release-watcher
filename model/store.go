package model

type Store interface {
	Exists(release Release) bool
	Mark(release Release) error
	Close()
}
