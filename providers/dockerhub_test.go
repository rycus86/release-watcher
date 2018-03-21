package providers

import (
	"io/ioutil"
	"net/http"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
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

	releases, err := provider.FetchReleases("rycus86", "grafana")
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
