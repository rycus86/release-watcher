package model

import (
	"fmt"
)

const (
	defaultFilterPatter = "^[0-9]+\\.[0-9]+\\.[0-9]+$"
)

type Project struct {
	Owner  string
	Repo   string
	Filter string
}

type Configuration struct {
	Releases map[string][]Project
	Path     string
}

func (p Project) String() string {
	if p.Owner != "" {
		return fmt.Sprintf("%s/%s", p.Owner, p.Repo)
	} else {
		return p.Repo
	}
}

func (p Project) GetFilter() string {
	if p.Filter != "" {
		return p.Filter
	} else {
		return defaultFilterPatter
	}
}
