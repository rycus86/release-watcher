package providers

import (
	"strings"

	"github.com/rycus86/release-watcher/model"
)

var providers []model.Provider

func GetProviders() []model.Provider {
	return providers
}

func GetProvider(name string) model.Provider {
	for _, provider := range providers {
		if strings.ToLower(provider.GetName()) == strings.ToLower(name) {
			return provider
		}
	}

	return nil
}

func RegisterProvider(provider model.Provider) {
	providers = append(providers, provider)
}

func InitializeProviders() {
	(&GitHubProvider{}).Initialize()
	(&DockerHubProvider{}).Initialize()
	(&PyPIProvider{}).Initialize()
	(&JetBrainsProvider{}).Initialize()
}
