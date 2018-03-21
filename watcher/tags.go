package watcher

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
)

type TagWatcher interface {
	FetchTags(project config.Project) ([]model.Tag, error)
}
