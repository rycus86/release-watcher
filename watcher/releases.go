package watcher

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
)

type ReleaseWatcher interface {
	FetchReleases(project config.Project) ([]model.Release, error)
}
