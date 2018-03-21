package config

import "testing"

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig("../testdata/sample_config.yml")
	if err != nil {
		t.Error("Failed to parse the configuration")
	}

	// Releases
	for provider, projects := range config.Releases {
		if provider == "github" {
			if len(projects) != 1 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].Owner != "docker" || projects[0].Repo != "docker-py" {
				t.Error("Unexpected project")
			}

		} else if provider == "dockerhub" {
			if len(projects) != 2 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].Owner != "" || projects[0].Repo != "nginx" {
				t.Error("Unexpected project")
			}

			if projects[1].Owner != "rycus86" || projects[1].Repo != "grafana" {
				t.Error("Unexpected project")
			}

		} else if provider == "pypi" {
			if len(projects) != 1 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].Owner != "" || projects[0].Repo != "flask" {
				t.Error("Unexpected project")
			}

		} else {
			t.Errorf("Unexpected provider: %s", provider)

		}
	}

	// Tags
	for provider, projects := range config.Tags {
		if provider == "github" {
			if len(projects) != 2 {
				t.Error("Unexpected number of projects")
			}

			if projects[0].Owner != "rycus86" || projects[0].Repo != "prometheus-flask-exporter" {
				t.Error("Unexpected project")
			}

			if projects[1].Owner != "rycus86" || projects[1].Repo != "ghost-client" {
				t.Error("Unexpected project")
			}

		} else {
			t.Errorf("Unexpected provider: %s", provider)

		}
	}
}
