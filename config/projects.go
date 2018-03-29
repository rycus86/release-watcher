package config

import (
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/providers"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func ParseConfigurationFile(path string) (*model.Configuration, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var perProvider map[string]interface{}
	err = yaml.Unmarshal(contents, &perProvider)
	if err != nil {
		return nil, err
	}

	configuration := model.Configuration{
		Releases: map[string][]model.GenericProject{},
	}
	configuration.Path = path

	for name, settings := range perProvider["releases"].(map[interface{}]interface{}) {
		providerName := name.(string)
		provider := providers.GetProvider(providerName)

		if provider == nil {
			log.Panicln("Invalid provider:", providerName)
		}

		for _, item := range settings.([]interface{}) {
			parsed := provider.Parse(item)
			if parsed == nil {
				log.Println("Failed to parse", item, "for", providerName)
				continue
			}

			configuration.Releases[providerName] = append(configuration.Releases[providerName], parsed)
		}
	}

	return &configuration, nil
}

func Reload(c *model.Configuration) error {
	if newConfig, err := ParseConfigurationFile(c.Path); err != nil {
		return err
	} else {
		c.Releases = newConfig.Releases
		return nil
	}
}
