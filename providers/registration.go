package providers

// TODO models?
type Provider interface {
	Initialize()
	GetName() string
}

var providers []Provider

// TODO is this unused?
func GetProviders() []Provider {
	return providers
}

func GetProvider(name string) Provider {
	for _, provider := range providers {
		if provider.GetName() == name {
			return provider
		}
	}

	return nil
}

func RegisterProvider(provider Provider) {
	providers = append(providers, provider)
}

func InitializeProviders() {
	(&GitHubProvider{}).Initialize()
	(&DockerHubProvider{}).Initialize()
	(&PyPIProvider{}).Initialize()
}
