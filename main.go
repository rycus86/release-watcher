package main

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/providers"
	"github.com/rycus86/release-watcher/watcher"
	"github.com/rycus86/release-watcher/store"
)

func main() {
	store.DbTest()

	providers.InitializeProviders()

	for _, pr := range providers.GetProviders() {
		println(pr.GetName())

		if releaseWatcher, ok := pr.(watcher.ReleaseWatcher); ok {
			releases, err := releaseWatcher.FetchReleases(config.Project{Owner: "rycus86", Repo: "grafana"})

			if err != nil {
				println("error:", err.Error())
			} else {
				for _, r := range releases {
					println(r.Name, r.Date.String(), r.URL)
				}
			}
		}

		if tagWatcher, ok := pr.(watcher.TagWatcher); ok {
			tags, err := tagWatcher.FetchTags(config.Project{Owner: "rycus86", Repo: "ghost-client"})

			if err != nil {
				println("error:", err.Error())
			} else {
				for _, t := range tags {
					println(t.Name, t.Date.String(), t.URL, t.Message)
				}
			}
		}
	}
}
