package watcher

import (
	"errors"
	"github.com/rycus86/release-watcher/model"
	"os"
	"testing"
	"time"
)

func TestWatchOnce(t *testing.T) {
	w := &mockWatcher{}

	project := &mockProject{}

	mr := []model.Release{
		{
			Provider: w,
			Project:  project,

			Name: "0.0.1",
			Date: time.Now().Add(-10 * time.Minute),
			URL:  "http://test.release/0.0.1",
		},
		{
			Provider: w,
			Project:  project,

			Name: "0.0.2",
			Date: time.Now().Add(-3 * time.Minute),
			URL:  "http://test.release/0.0.2",
		},
	}

	w.Releases = mr

	out := make(chan []model.Release, 1)
	done := make(chan struct{})

	go WatchReleases(w, project, out, done)

	close(done)
	defer close(out)

	releases := <-out
	if len(releases) != 2 {
		t.Error("Invalid releases found:", releases)
	}

	for _, release := range releases {
		if release.Provider.GetName() != "Mock" {
			t.Error("Invalid provider for release:", release.Provider.GetName())
		}

		if release.Project.String() != project.String() {
			t.Error("Invalid project:", release.Project)
		}
	}

	if releases[0].Name != "0.0.2" {
		t.Error("Invalid release version:", releases[0].Name)
	}
	if releases[1].Name != "0.0.1" {
		t.Error("Invalid release version:", releases[1].Name)
	}

	if releases[0].URL != "http://test.release/0.0.2" {
		t.Error("Invalid release version:", releases[0].URL)
	}
	if releases[1].URL != "http://test.release/0.0.1" {
		t.Error("Invalid release version:", releases[1].URL)
	}
}

func TestWatchTicker(t *testing.T) {
	os.Setenv("CHECK_INTERVAL", "1ms")
	defer os.Unsetenv("CHECK_INTERVAL")

	w := &mockWatcher{}

	project := &mockProject{}

	out := make(chan []model.Release, 1)
	done := make(chan struct{})

	go WatchReleases(w, project, out, done)

	defer close(done)
	defer close(out)

	w.NextReleases = []model.Release{
		{
			Provider: w,
			Project:  project,

			Name: "1.2.3",
			URL:  "http://test.ticker/1.2.3",
		},
	}

	releases := <-out
	if len(releases) != 1 {
		t.Error("Invalid releases found:", releases)
	}

	if releases[0].Name != "1.2.3" {
		t.Error("Unexpected release:", releases[0].Name)
	}

	if w.FetchCount != 2 {
		t.Error("Unexpected number of fetches:", w.FetchCount)
	}
}

func TestNoReleasesOnErrors(t *testing.T) {
	w := &mockWatcher{
		Error: errors.New("Test Error"),
	}

	project := &mockProject{}

	out := make(chan []model.Release, 1)
	done := make(chan struct{})

	close(done)
	defer close(out)

	WatchReleases(w, project, out, done)

	if w.ErrorCount != 1 {
		t.Error("Expected to return an error")
	}
}

type mockWatcher struct {
	Releases     []model.Release
	NextReleases []model.Release
	Error        error
	FetchCount   int
	ErrorCount   int
}

func (m *mockWatcher) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	m.FetchCount++
	if m.Error != nil {
		m.ErrorCount++
	}

	defer func() {
		m.Releases = m.NextReleases
	}()

	return m.Releases, m.Error
}

func (m *mockWatcher) GetName() string {
	return "Mock"
}

func (m *mockWatcher) Initialize() {}

func (m *mockWatcher) Parse(interface{}) model.GenericProject {
	return nil
}

type mockProject struct{}

func (p *mockProject) String() string {
	return "mock/repo"
}

func (p *mockProject) GetFilter() string {
	return ".*"
}
