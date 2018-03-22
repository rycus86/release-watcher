package providers

import (
	"github.com/rycus86/release-watcher/model"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetchDockerHubReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/dockerhub_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://hub.docker.com/v2/repositories/rycus86/grafana/tags/",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := DockerHubProvider{
		client: &http.Client{},
	}

	releases, err := provider.FetchReleases(model.Project{Owner: "rycus86", Repo: "grafana"})
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}

	if len(releases) != 18 {
		t.Error("Wrong number of results")
	}

	sample := releases[3]

	if sample.Name != "5.0.2" {
		t.Errorf("Unexpected name: %s", sample.Name)
	}

	if sample.Date.Year() != 2018 || sample.Date.Month() != 3 || sample.Date.Day() != 14 {
		t.Errorf("Unexpected date: %s", sample.Date.String())
	}

	if sample.URL != "https://hub.docker.com/r/rycus86/grafana/tags/" {
		t.Errorf("Unexpected URL: %s", sample.URL)
	}
}

func TestFetchForLibraryImage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provider := DockerHubProvider{
		client: &http.Client{},
	}

	httpmock.RegisterResponder(
		"GET", "https://hub.docker.com/v2/repositories/_/nginx/tags/",
		httpmock.NewStringResponder(200, "{}"),
	)

	_, err := provider.FetchReleases(model.Project{Owner: "_", Repo: "nginx"})
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}
}
