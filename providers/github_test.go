package providers

import (
	"io/ioutil"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"testing"
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

	releases, err := provider.FetchReleases("docker", "docker-py")
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

func TestFetchGitHubTags(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/github_tags.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://api.github.com/repos/docker/docker-py/tags",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := GitHubProvider{}
	provider.Initialize()

	releases, err := provider.FetchTags("docker", "docker-py")
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}

	if len(releases) != 30 {
		t.Error("Wrong number of results")
	}

	sample := releases[1]

	if sample.Name != "3.1.2" {
		t.Errorf("Unexpected name: %s", sample.Name)
	}

	if sample.URL != "https://github.com/docker/docker-py/commit/88b0d697aa5386c2ef90a5b480cd400ce5a32646" {
		t.Errorf("Unexpected URL: %s", sample.URL)
	}
}
