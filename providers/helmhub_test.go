package providers

import (
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestFetchHelmHubReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/helmhub_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://hub.helm.sh/api/chartsvc/v1/charts/argo/argo/versions",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := HelmHubProvider{
		client: &http.Client{},
	}

	releases, err := provider.FetchReleases(&HelmHubProject{Repo: "argo", Chart: "argo"})
	if err != nil {
		t.Error("Failed:", err)
	}

	if len(releases) != 22 {
		t.Error("Unexpected number of releases:", len(releases))
	}

	if releases[0].Name != "0.7.2" {
		t.Error("Unexpected release:", releases[0].Name)
	}
	if releases[3].Name != "0.6.8" {
		t.Error("Unexpected release:", releases[1].Name)
	}

	if releases[0].Date.Year() != 2020 || releases[0].Date.Month() != 3 || releases[0].Date.Day() != 19 {
		t.Errorf("Unexpected date: %s", releases[0].Date.String())
	}

	if releases[0].URL != "https://hub.helm.sh/charts/argo/argo/0.7.2" {
		t.Error("Unexpected URL:", releases[0].URL)
	}
}
