package main

import (
	"github.com/rycus86/release-watcher/providers"
	"github.com/rycus86/release-watcher/watcher"
)

func main() {
	providers.InitializeProviders()

	for _, pr := range providers.GetProviders() {
		println(pr.GetName())

		if releaseWatcher, ok := pr.(watcher.ReleaseWatcher); ok {
			releases, err := releaseWatcher.FetchReleases("rycus86", "grafana")

			if err != nil {
				println("error:", err.Error())
			} else {
				for _, r := range releases {
					println(r.Name, r.Date.String(), r.URL)
				}
			}
		}

		if tagWatcher, ok := pr.(watcher.TagWatcher); ok {
			tags, err := tagWatcher.FetchTags("rycus86", "ghost-client")

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
