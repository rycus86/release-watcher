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

type HelmHubProvider struct {
	client *http.Client
}

type helmHubResponse struct {
	Data []struct {
		Attributes struct {
			Version string `json:"version"`
			Created string `json:"created"`
		} `json:"attributes"`
	} `json:"data"`
}

type HelmHubProject struct {
	model.BaseProject `mapstructure:",squash"`

	Repo  string
	Chart string
}

func (p HelmHubProject) String() string {
	if p.Repo != "" {
		return fmt.Sprintf("%s/%s", p.Repo, p.Chart)
	} else {
		return p.Chart
	}
}

func (provider *HelmHubProvider) Initialize() {
	provider.client = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/dockerhub"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}

	RegisterProvider(provider)
}

func (provider *HelmHubProvider) GetName() string {
	return "HelmHub"
}

func (provider *HelmHubProvider) Parse(input interface{}) model.GenericProject {
	var project HelmHubProject

	err := mapstructure.Decode(input, &project)
	if err != nil {
		return nil
	}

	return &project
}

func (provider *HelmHubProvider) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*HelmHubProject)

	apiURL := fmt.Sprintf(
		"https://hub.helm.sh/api/chartsvc/v1/charts/%s/%s/versions",
		project.Repo, project.Chart)

	response, err := provider.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse = helmHubResponse{}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	for _, release := range apiResponse.Data {
		url := fmt.Sprintf("https://hub.helm.sh/charts/%s/%s/%s", project.Repo, project.Chart, release.Attributes.Version)
		created, err := time.Parse(time.RFC3339Nano, release.Attributes.Created)
		if err != nil {
			created = time.Now()
		}

		releases = append(releases, model.Release{
			Name: release.Attributes.Version,
			URL:  url,
			Date: created,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
