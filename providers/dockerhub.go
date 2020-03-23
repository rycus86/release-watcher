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

type DockerHubProvider struct {
	client   *http.Client
	pageSize int
}

type dockerHubTagsResponse struct {
	Results []struct {
		Name        string `json:"name"`
		LastUpdated string `json:"last_updated"`
	} `json:"results"`
}

type DockerHubProject struct {
	model.BaseProject `mapstructure:",squash"`

	Owner string
	Repo  string
}

func (p DockerHubProject) String() string {
	if p.Owner != "" {
		return fmt.Sprintf("%s/%s", p.Owner, p.Repo)
	} else {
		return p.Repo
	}
}

func (provider *DockerHubProvider) Initialize() {
	provider.client = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/dockerhub"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}
	provider.pageSize = env.GetInt("PAGE_SIZE", "/var/secrets/dockerhub", 50)

	RegisterProvider(provider)
}

func (provider *DockerHubProvider) GetName() string {
	return "DockerHub"
}

func (provider *DockerHubProvider) Parse(input interface{}) model.GenericProject {
	var project DockerHubProject

	err := mapstructure.Decode(input, &project)
	if err != nil {
		return nil
	}

	return &project
}

func (provider *DockerHubProvider) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*DockerHubProject)

	owner := project.Owner
	if owner == "" {
		owner = "library"
	}

	apiUrl := fmt.Sprintf(
		"https://hub.docker.com/v2/repositories/%s/%s/tags/?page=1&page_size=%d",
		owner, project.Repo, provider.pageSize)

	response, err := provider.client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse = dockerHubTagsResponse{}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	for _, release := range apiResponse.Results {
		url := fmt.Sprintf("https://hub.docker.com/r/%s/%s/tags/", owner, project.Repo)
		published, err := time.Parse(time.RFC3339Nano, release.LastUpdated)
		if err != nil {
			published = time.Now()
		}

		releases = append(releases, model.Release{
			Name: release.Name,
			URL:  url,
			Date: published,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
