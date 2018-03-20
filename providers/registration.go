package providers

type Provider interface {
	Initialize()
	GetName() string
}

var providers []Provider

func GetProviders() []Provider {
	return providers
}

func RegisterProvider(provider Provider) {
	providers = append(providers, provider)
}

func InitializeProviders() {
	(&GitHubProvider{}).Initialize()
	(&DockerHubProvider{}).Initialize()
}
