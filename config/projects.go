package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Project struct {
	Owner string
	Repo  string
}

type Configuration struct {
	Releases map[string][]Project
	Tags     map[string][]Project
}

func ParseConfig(path string) (*Configuration, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configuration Configuration
	if err := yaml.Unmarshal(contents, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}
