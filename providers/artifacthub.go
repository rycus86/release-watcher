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

type ArtifactHubProvider struct {
	client *http.Client
}

type ArtifactHubResponse struct {
	AvailableVersions []struct {
		Version string `json:"version"`
		Created int64  `json:"ts"`
	} `json:"available_versions"`
}

type ArtifactHubProject struct {
	model.BaseProject `mapstructure:",squash"`

	Repo  string
	Chart string
}

func (p ArtifactHubProject) String() string {
	if p.Repo != "" {
		return fmt.Sprintf("%s/%s", p.Repo, p.Chart)
	} else {
		return p.Chart
	}
}

func (provider *ArtifactHubProvider) Initialize() {
	provider.client = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/dockerhub"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}

	RegisterProvider(provider)
}

func (provider *ArtifactHubProvider) GetName() string {
	return "ArtifactHub"
}

func (provider *ArtifactHubProvider) Parse(input interface{}) model.GenericProject {
	var project ArtifactHubProject

	err := mapstructure.Decode(input, &project)
	if err != nil {
		return nil
	}

	return &project
}

func (provider *ArtifactHubProvider) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*ArtifactHubProject)

	apiURL := fmt.Sprintf(
		"https://artifacthub.io/api/v1/packages/helm/%s/%s",
		project.Repo, project.Chart)

	response, err := provider.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse = ArtifactHubResponse{}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	for _, release := range apiResponse.AvailableVersions {
		url := fmt.Sprintf("https://artifacthub.io/packages/helm/%s/%s/%s", project.Repo, project.Chart, release.Version)
		created := time.Unix(release.Created, 0)
		if err != nil {
			created = time.Now()
		}

		releases = append(releases, model.Release{
			Name: release.Version,
			URL:  url,
			Date: created,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
