package store

import (
	"testing"
	"github.com/rycus86/release-watcher/model"
)

func TestInitialize(t *testing.T) {
	db, err := InitForTesting()
	if err != nil {
		t.Error("Failed to initialize the store")
	}
	defer db.Close()
}

func TestGetSet(t *testing.T) {
	db, err := InitForTesting()
	if err != nil {
		t.Error("Failed to initialize the store:", err)
	}
	defer db.Close()

	err = db.Set("key:x", "value:x")
	if err != nil {
		t.Error("Failed to set value:", err)
	}

	err = db.Set("key:y", "value:y")
	if err != nil {
		t.Error("Failed to set value:", err)
	}

	value := db.Get("key:x")
	if value != "value:x" {
		t.Error("Unexpected value:", value)
	}
}

func InitForTesting() (model.Store, error) {
	return Initialize("file::memory:?cache=shared")
}
