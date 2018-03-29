package model

const (
	defaultFilterPatter = "^[0-9]+\\.[0-9]+\\.[0-9]+$"
)

type GenericProject interface {
	String() string
	GetFilter() string
}

type Project struct {
	Name   string
	Filter string
}

type Configuration struct {
	Releases map[string][]GenericProject
	Path     string
}

func (p Project) String() string {
	return p.Name
}

func (p Project) GetFilter() string {
	if p.Filter != "" {
		return p.Filter
	} else {
		return defaultFilterPatter
	}
}
