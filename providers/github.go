package providers

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
)

type GitHubProvider struct {
	client *github.Client
}

func (provider *GitHubProvider) Initialize() {
	username := config.Lookup("GITHUB_USERNAME", "/var/secrets/github", "")
	password := config.Lookup("GITHUB_PASSWORD", "/var/secrets/github", "")

	if username != "" && password != "" {
		transport := github.BasicAuthTransport{
			Username: "x",
			Password: "x",
		}

		provider.client = github.NewClient(transport.Client())

	} else {
		provider.client = github.NewClient(nil)

	}

	RegisterProvider(provider)
}

func (provider *GitHubProvider) GetName() string {
	return "GitHub"
}

func (provider *GitHubProvider) FetchReleases(project model.Project) ([]model.Release, error) {
	var releases []model.Release

	ctx, cancel := context.WithTimeout(
		context.Background(), config.GetTimeout("HTTP_TIMEOUT", "/var/secrets/github"),
	)
	defer cancel()

	ghReleases, _, err := provider.client.Repositories.ListReleases(ctx, project.Owner, project.Repo, &github.ListOptions{PerPage: 50})
	if err != nil {
		return nil, err
	}

	for _, release := range ghReleases {
		releases = append(releases, model.Release{
			Name: release.GetName(),
			URL:  release.GetHTMLURL(),
			Date: release.GetPublishedAt().Time,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
