package model

const (
	DefaultFilterPattern = "^[0-9]+\\.[0-9]+\\.[0-9]+$"
)

type GenericProject interface {
	String() string
	GetFilter() string
}

type Configuration struct {
	Releases map[string][]GenericProject
	Path     string
}
