package providers

import (
	"context"
	"fmt"
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
	return "github"
}

func (provider *GitHubProvider) FetchReleases(project config.Project) ([]model.Release, error) {
	var releases []model.Release

	ctx, cancel := context.WithTimeout(
		context.Background(), config.GetDuration("HTTP_TIMEOUT", "/var/secrets/github"),
	)
	defer cancel()

	ghReleases, _, err := provider.client.Repositories.ListReleases(ctx, project.Owner, project.Repo, nil)
	if err != nil {
		return nil, err
	}

	for _, release := range ghReleases {
		releases = append(releases, model.Release{
			Name: release.GetName(),
			URL:  release.GetHTMLURL(),
			Date: release.GetPublishedAt().Time,
		})
	}

	return releases, nil
}

func (provider *GitHubProvider) FetchTags(project config.Project) ([]model.Tag, error) {
	var tags []model.Tag

	ctx, cancel := context.WithTimeout(
		context.Background(), config.GetDuration("HTTP_TIMEOUT", "/var/secrets/github"),
	)
	defer cancel()

	ghTags, _, err := provider.client.Repositories.ListTags(ctx, project.Owner, project.Repo, nil)
	if err != nil {
		return nil, err
	}

	for _, tag := range ghTags {
		url := tag.GetCommit().GetHTMLURL()
		if url == "" {
			url = fmt.Sprintf("https://github.com/%s/%s/commit/%s", project.Owner, project.Repo, tag.GetCommit().GetSHA())
		}

		tags = append(tags, model.Tag{
			Name:    tag.GetName(),
			URL:     url,
			Message: tag.GetCommit().GetMessage(),
		})
	}

	return tags, nil
}
