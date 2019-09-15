package store

import (
	"github.com/rycus86/release-watcher/model"
	"testing"
)

func TestInitialize(t *testing.T) {
	db, err := InitForTesting()
	if err != nil {
		t.Error("Failed to initialize the store")
	}
	defer db.Close()
}

func TestExistsAndMark(t *testing.T) {
	db, err := InitForTesting()
	if err != nil {
		t.Error("Failed to initialize the store:", err)
	}
	defer db.Close()

	release1 := model.Release{
		Provider: mockProvider{
			Name: "TestProvider",
		},
		Project: &mockProject{
			Name: "sample/repo",
		},
		Name: "test-tag",
	}

	release2 := model.Release{
		Provider: mockProvider{
			Name: "TestProvider",
		},
		Project: &mockProject{
			Name: "sample/alt",
		},
		Name: "1.0.0",
	}

	err = db.Mark(release1)
	if err != nil {
		t.Error("Failed to mark a release:", err)
	}

	err = db.Mark(release2)
	if err != nil {
		t.Error("Failed to mark a release:", err)
	}

	if !db.Exists(release1) || !db.Exists(release2) {
		t.Error("Saved release not found")
	}
}

func InitForTesting() (model.Store, error) {
	return Initialize("file::memory:?cache=shared")
}

type mockProvider struct {
	Name string
}

func (p mockProvider) Initialize() {

}

func (p mockProvider) GetName() string {
	return p.Name
}

func (p mockProvider) Parse(interface{}) model.GenericProject {
	return nil
}

type mockProject struct {
	Name string
}

func (p *mockProject) String() string {
	return p.Name
}

func (p *mockProject) GetFilter() string {
	return ".*"
}

func (p *mockProject) GetWebhooks() []string {
	return nil
}
