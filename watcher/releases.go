package watcher

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"log"
	"sort"
	"time"
)

type ReleaseWatcher interface {
	FetchReleases(project model.Project) ([]model.Release, error)
}

func WatchReleases(w ReleaseWatcher, project model.Project, outChannel chan<- []model.Release, done <-chan struct{}) {
	fetchNow(w, project, outChannel)

	ticker := time.NewTicker(config.GetInterval("CHECK_INTERVAL", "/var/secrets/release-watcher"))

	for {
		select {
		case <-ticker.C:
			fetchNow(w, project, outChannel)

		case <-done:
			ticker.Stop()
			return

		}
	}
}

func fetchNow(w ReleaseWatcher, project model.Project, outChannel chan<- []model.Release) {
	log.Println("Fetching releases for", project, "using", w.(model.Provider).GetName())

	releases, err := w.FetchReleases(project)
	if err != nil {
		log.Println("Failed to fetch the releases of", project, ":", err)
		return
	}

	SortReleases(Releases(releases))

	outChannel <- releases
}

type Releases []model.Release

func SortReleases(releases []model.Release) {
	sort.Stable(sort.Reverse(Releases(releases)))
}

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Swap(a, b int) {
	r[a], r[b] = r[b], r[a]
}

func (r Releases) Less(a, b int) bool {
	r1, r2 := r[a], r[b]

	if r1.Date.UnixNano() != r2.Date.UnixNano() {
		return r1.Date.UnixNano() < r2.Date.UnixNano()
	} else {
		return r1.Name < r2.Name
	}
}
