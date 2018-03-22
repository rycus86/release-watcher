package providers

import (
	"encoding/json"
	"fmt"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"net/http"
	"time"
)

const configPath = "/var/secrets/dockerhub"

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

func (provider *DockerHubProvider) Initialize() {
	provider.client = &http.Client{
		Timeout: config.GetTimeout("HTTP_TIMEOUT", configPath),
	}
	provider.pageSize = config.GetInt("PAGE_SIZE", configPath, 50)

	RegisterProvider(provider)
}

func (provider *DockerHubProvider) GetName() string {
	return "DockerHub"
}

func (provider *DockerHubProvider) FetchReleases(project model.Project) ([]model.Release, error) {
	var releases []model.Release

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
