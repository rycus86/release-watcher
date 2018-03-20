package watcher

import "github.com/rycus86/release-watcher/model"

type ReleaseWatcher interface {
	FetchReleases(owner string, repo string) ([]model.Release, error)
}
