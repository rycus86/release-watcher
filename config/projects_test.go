package config

import (
	"testing"

	"github.com/rycus86/release-watcher/providers"
)

func TestParseConfig(t *testing.T) {
	providers.InitializeProviders()

	configuration, err := ParseConfigurationFile("../testdata/sample_config.yml")
	if err != nil {
		t.Error("Failed to parse the configuration")
	}

	if len(configuration.Releases) != 4 {
		t.Error("Unexpected number of providers:", len(configuration.Releases))
	}

	// Releases
	for provider, projects := range configuration.Releases {
		if provider == "github" {
			if len(projects) != 2 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].String() != "docker/docker-py" {
				t.Error("Unexpected project:", projects[0])
			}

			if filter := projects[0].GetFilter(); filter != "[2-9]+\\.[0-9]+\\.[0-9]+" {
				t.Error("Unexpected filter:", filter)
			}

			if wh := projects[0].GetWebhooks(); len(wh) != 1 || wh[0] != "http://example.local/webhooks/test" {
				t.Error("Unexpected webhooks configuration:", wh)
			}

			if projects[1].String() != "nginx/nginx" {
				t.Error("Unexpected project:", projects[0])
			}

		} else if provider == "dockerhub" {
			if len(projects) != 2 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].String() != "nginx" {
				t.Error("Unexpected project:", projects[0])
			}

			if projects[1].String() != "rycus86/grafana" {
				t.Error("Unexpected project:", projects[1])
			}

		} else if provider == "pypi" {
			if len(projects) != 1 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].String() != "flask" {
				t.Error("Unexpected project:", projects[0])
			}

		} else if provider == "jetbrains" {
			if len(projects) != 3 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].String() != "GoLand" {
				t.Error("Unexpected project:", projects[0])
			}
			if projects[1].String() != "IntelliJ IDEA" {
				t.Error("Unexpected project:", projects[1])
			}
			if projects[2].String() != "PyCharm" {
				t.Error("Unexpected project:", projects[2])
			}

		} else {
			t.Errorf("Unexpected provider: %s", provider)

		}
	}
}
