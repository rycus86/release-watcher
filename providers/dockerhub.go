package providers

import (
	"encoding/json"
	"fmt"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"net/http"
	"time"
)

type DockerHubProvider struct {
	client *http.Client
}

type dockerHubTagsResponse struct {
	Results []struct {
		Name        string `json:"name"`
		LastUpdated string `json:"last_updated"`
	} `json:"results"`
}

func (provider *DockerHubProvider) Initialize() {
	provider.client = &http.Client{
		Timeout: config.GetDuration("HTTP_TIMEOUT", "/var/secrets/dockerhub"),
	}

	RegisterProvider(provider)
}

func (provider *DockerHubProvider) GetName() string {
	return "dockerhub"
}

func (provider *DockerHubProvider) FetchReleases(project config.Project) ([]model.Release, error) {
	var releases []model.Release

	owner := project.Owner
	if owner == "" {
		owner = "_"
	}

	apiUrl := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/", project.Owner, project.Repo)
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
		url := fmt.Sprintf("https://hub.docker.com/r/%s/%s/tags/", project.Owner, project.Repo)
		published, err := time.Parse(time.RFC3339Nano, release.LastUpdated)
		if err != nil {
			published = time.Now()
		}

		releases = append(releases, model.Release{
			Name: release.Name,
			URL:  url,
			Date: published,
		})
	}

	return releases, nil
}
