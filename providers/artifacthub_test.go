package providers

import (
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestFetchArtifactHubReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/artifacthub_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://artifacthub.io/api/v1/packages/helm/argo/argo",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := ArtifactHubProvider{
		client: &http.Client{},
	}

	releases, err := provider.FetchReleases(&ArtifactHubProject{Repo: "argo", Chart: "argo"})
	if err != nil {
		t.Error("Failed:", err)
	}

	if len(releases) != 80 {
		t.Error("Unexpected number of releases:", len(releases))
	}

	if releases[0].Name != "0.16.8" {
		t.Error("Unexpected release:", releases[0].Name)
	}
	if releases[3].Name != "0.16.7" {
		t.Error("Unexpected release:", releases[3].Name)
	}

	if releases[0].Date.Year() != 2021 || releases[0].Date.Month() != 4 || releases[0].Date.Day() != 22 {
		t.Errorf("Unexpected date: %s", releases[0].Date.String())
	}

	if releases[0].URL != "https://artifacthub.io/packages/helm/argo/argo/0.16.8" {
		t.Error("Unexpected URL:", releases[0].URL)
	}
}
