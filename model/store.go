package model

type Store interface {
	Get(key string) string
	Set(key string, value string) error
	Close()
}
