package config

import (
	"github.com/rycus86/release-watcher/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ParseConfigurationFile(path string) (*model.Configuration, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configuration model.Configuration
	if err := yaml.Unmarshal(contents, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}
