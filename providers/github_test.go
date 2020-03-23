package providers

import (
	"io/ioutil"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestFetchGitHubReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/github_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://api.github.com/repos/docker/docker-py/releases",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := GitHubProvider{}
	provider.Initialize()

	releases, err := provider.FetchReleases(&GitHubProject{Owner: "docker", Repo: "docker-py"})
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}

	if len(releases) != 30 {
		t.Error("Wrong number of results")
	}

	sample := releases[1]

	if sample.Name != "3.1.1" {
		t.Errorf("Unexpected name: %s", sample.Name)
	}

	if sample.Date.Year() != 2018 || sample.Date.Month() != 3 || sample.Date.Day() != 5 {
		t.Errorf("Unexpected date: %s", sample.Date.String())
	}

	if sample.URL != "https://github.com/docker/docker-py/releases/tag/3.1.1" {
		t.Errorf("Unexpected URL: %s", sample.URL)
	}
}

func TestFetchGitHubReleasesWithTagName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/github_traefik_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://api.github.com/repos/containous/traefik/releases",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := GitHubProvider{}
	provider.Initialize()

	releases, err := provider.FetchReleases(&GitHubProject{Owner: "containous", Repo: "traefik"})
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}

	if len(releases) != 30 {
		t.Error("Wrong number of results")
	}

	sample := releases[0]

	if sample.Name != "v1.6.4" {
		t.Errorf("Unexpected name: %s", sample.Name)
	}

	if sample.Date.Year() != 2018 || sample.Date.Month() != 6 || sample.Date.Day() != 15 {
		t.Errorf("Unexpected date: %s", sample.Date.String())
	}

	if sample.URL != "https://github.com/containous/traefik/releases/tag/v1.6.4" {
		t.Errorf("Unexpected URL: %s", sample.URL)
	}
}
