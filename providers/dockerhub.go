package providers

import (
	"encoding/json"
	"fmt"
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
		Timeout: 30 * time.Second,
	}

	RegisterProvider(provider)
}

func (provider *DockerHubProvider) GetName() string {
	return "dockerhub"
}

func (provider *DockerHubProvider) FetchReleases(owner string, repo string) ([]model.Release, error) {
	var releases []model.Release

	apiUrl := provider.getUrl(fmt.Sprintf("/v2/repositories/%s/%s/tags/", owner, repo))
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
		url := provider.getUrl(fmt.Sprintf("/r/%s/%s/tags/", owner, repo))
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

func (provider *DockerHubProvider) getUrl(path string) string {
	return fmt.Sprintf("https://hub.docker.com%s", path)
}
