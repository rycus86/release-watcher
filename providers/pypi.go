package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
)

type PyPIProvider struct {
	client *http.Client
}

type pypiResponse struct {
	Releases map[string][]struct {
		UploadTime string `json:"upload_time"`
	} `json:"releases"`
}

type PyPIProject struct {
	model.BaseProject `mapstructure:",squash"`

	Name string
}

func (p PyPIProject) String() string {
	return p.Name
}

func (provider *PyPIProvider) Initialize() {
	provider.client = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/pypi"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}

	RegisterProvider(provider)
}

func (provider *PyPIProvider) GetName() string {
	return "PyPI"
}

func (provider *PyPIProvider) Parse(input interface{}) model.GenericProject {
	var project PyPIProject

	err := mapstructure.Decode(input, &project)
	if err != nil {
		return nil
	}

	return &project
}

func (provider *PyPIProvider) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*PyPIProject)

	apiUrl := fmt.Sprintf("https://pypi.python.org/pypi/%s/json", project.Name)
	response, err := provider.client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse = pypiResponse{}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	for version, items := range apiResponse.Releases {
		if len(items) == 0 {
			continue
		}

		release := items[0]

		published, err := time.Parse("2006-01-02T15:04:05", release.UploadTime)
		if err != nil {
			published = time.Now()
		}

		releases = append(releases, model.Release{
			Name: version,
			URL:  fmt.Sprintf("https://pypi.python.org/pypi/%s/%s", project.Name, version),
			Date: published,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
