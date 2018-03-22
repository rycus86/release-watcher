package model

import (
	"fmt"
)

type Project struct {
	Owner string
	Repo  string
}

type Configuration struct {
	Releases map[string][]Project
}

func (p Project) String() string {
	if p.Owner != "" {
		return fmt.Sprintf("%s/%s", p.Owner, p.Repo)
	} else {
		return p.Repo
	}
}
