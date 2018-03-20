package watcher

import "github.com/rycus86/release-watcher/model"

type TagWatcher interface {
	FetchTags(owner string, repo string) ([]model.Tag, error)
}
