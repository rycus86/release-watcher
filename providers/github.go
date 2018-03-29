package providers

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/mapstructure"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
)

type GitHubProvider struct {
	client *github.Client
}

type GitHubProject struct {
	model.Project

	Owner string
	Repo  string
}

func (p GitHubProject) String() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Repo)
}

func (provider *GitHubProvider) Initialize() {
	username := env.Lookup("GITHUB_USERNAME", "/var/secrets/github", "")
	password := env.Lookup("GITHUB_PASSWORD", "/var/secrets/github", "")

	if username != "" && password != "" {
		transport := github.BasicAuthTransport{
			Username: username,
			Password: password,
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

func (provider *GitHubProvider) Parse(input interface{}) model.GenericProject {
	var project GitHubProject

	err := mapstructure.Decode(input, &project)
	if err != nil {
		return nil
	}

	return &project
}

func (provider *GitHubProvider) FetchReleases(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*GitHubProject)

	ctx, cancel := context.WithTimeout(
		context.Background(), env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/github"),
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
