package providers

import (
	"github.com/rycus86/release-watcher/model"
	"testing"
)

func TestRegisterProvider(t *testing.T) {
	providers = make([]model.Provider, 0)

	RegisterProvider(MockProvider{Name: "Test1"})

	if len(providers) != 1 {
		t.Error("Failed to register provider")
	}

	if providers[0].GetName() != "Test1" {
		t.Error("Invalid provider name")
	}

	RegisterProvider(MockProvider{Name: "Test2"})

	if len(providers) != 2 {
		t.Error("Failed to register provider")
	}

	if providers[1].GetName() != "Test2" {
		t.Error("Invalid provider name")
	}
	if providers[0].GetName() != "Test1" {
		t.Error("Invalid provider name")
	}
}

func TestGetProviders(t *testing.T) {
	providers = make([]model.Provider, 0)

	RegisterProvider(MockProvider{Name: "Test1"})
	RegisterProvider(MockProvider{Name: "Test2"})
	RegisterProvider(MockProvider{Name: "Test3"})

	registered := GetProviders()

	if len(registered) != 3 {
		t.Error("Unexpected providers")
	}

	if registered[0].GetName() != "Test1" {
		t.Error("Invalid provider name")
	}
	if registered[1].GetName() != "Test2" {
		t.Error("Invalid provider name")
	}
	if registered[2].GetName() != "Test3" {
		t.Error("Invalid provider name")
	}
}

func TestGetProvider(t *testing.T) {
	providers = make([]model.Provider, 0)

	RegisterProvider(MockProvider{Name: "Test1"})
	RegisterProvider(MockProvider{Name: "Test2"})

	p := GetProvider("Test1")
	if p == nil {
		t.Error("Provider not found")
	}
	if p.GetName() != "Test1" {
		t.Error("Invalid provider found")
	}

	p = GetProvider("Test2")
	if p == nil {
		t.Error("Provider not found")
	}
	if p.GetName() != "Test2" {
		t.Error("Invalid provider found")
	}

	p = GetProvider("NotFound")
	if p != nil {
		t.Error("Unexpected provider found")
	}
}

type MockProvider struct {
	Name        string
	Initialized bool
}

func (p MockProvider) Initialize() {
	p.Initialized = true
}

func (p MockProvider) GetName() string {
	return p.Name
}

func (p MockProvider) Parse(interface{}) model.GenericProject {
	return nil
}
