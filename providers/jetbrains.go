package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"net/http"
	"strings"
	"time"
)

type JetBrainsProvider struct {
	client *http.Client
}

type jetBrainsResponse struct {
	Releases []struct {
		Version   string `json:"version"`
		Build     string `json:"build"`
		Date      string `json:"date"`
		Type      string `json:"type"`
		Downloads struct {
			Linux struct {
				Link string `json:"link"`
			} `json:"linux"`
		} `json:"downloads"`
	} `json:"releases"`
}

func (provider *JetBrainsProvider) Initialize() {
	provider.client = &http.Client{
		Timeout: config.GetTimeout("HTTP_TIMEOUT", configPath),
	}

	RegisterProvider(provider)
}

func (provider *JetBrainsProvider) GetName() string {
	return "JetBrains"
}

func (provider *JetBrainsProvider) FetchReleases(project model.Project) ([]model.Release, error) {
	var releases []model.Release

	apiUrl := fmt.Sprintf(
		"https://data.services.jetbrains.com/products?code=%s",
		strings.ToUpper(project.Repo),
	)

	response, err := provider.client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse = make([]jetBrainsResponse, 1)
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse) != 1 {
		return nil, errors.New(fmt.Sprintf("unexpected number of response objects: %d", len(apiResponse)))
	}

	for _, release := range apiResponse[0].Releases {
		published, err := time.Parse("2006-01-02", release.Date)
		if err != nil {
			published = time.Now()
		}

		info := release.Build
		if release.Type != "release" {
			info = fmt.Sprintf("%s %s", info, release.Type)
		}

		releases = append(releases, model.Release{
			Name: fmt.Sprintf("%s (%s)", release.Version, info),
			URL:  release.Downloads.Linux.Link,
			Date: published,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
