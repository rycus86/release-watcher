package providers

import (
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/watcher"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetchPyPIReleases(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testdata, err := ioutil.ReadFile("../testdata/pypi_releases.json")
	if err != nil {
		t.Errorf("Failed to load test data: %s", err)
	}

	httpmock.RegisterResponder(
		"GET", "https://pypi.python.org/pypi/prometheus-flask-exporter/json",
		httpmock.NewStringResponder(200, string(testdata)),
	)

	provider := PyPIProvider{
		client: &http.Client{},
	}

	releases, err := provider.FetchReleases(&model.Project{Name: "prometheus-flask-exporter"})
	if err != nil {
		t.Errorf("Failed to fetch releases: %s", err)
	}

	if len(releases) != 13 {
		t.Error("Wrong number of results")
	}

	watcher.SortReleases(releases)

	sample := releases[1]

	if sample.Name != "0.2.0" {
		t.Errorf("Unexpected name: %s", sample.Name)
	}

	if sample.Date.Year() != 2018 || sample.Date.Month() != 2 || sample.Date.Day() != 27 {
		t.Errorf("Unexpected date: %s", sample.Date.String())
	}

	if sample.URL != "https://pypi.python.org/pypi/prometheus-flask-exporter/0.2.0" {
		t.Errorf("Unexpected URL: %s", sample.URL)
	}
}
