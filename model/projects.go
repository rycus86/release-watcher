package model

const (
	DefaultFilterPattern = "^[0-9]+\\.[0-9]+\\.[0-9]+$"
)

type GenericProject interface {
	String() string
	GetFilter() string
	GetWebhooks() []string
}

type BaseProject struct {
	Filter   string
	Webhooks []string
}

func (p BaseProject) GetFilter() string {
	if p.Filter != "" {
		return p.Filter
	}

	return DefaultFilterPattern
}

func (p BaseProject) GetWebhooks() []string {
	return p.Webhooks
}

type Configuration struct {
	Releases map[string][]GenericProject
	Path     string
}
