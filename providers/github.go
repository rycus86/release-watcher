package providers

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/rycus86/release-watcher/model"
)

type GitHubProvider struct {
	client *github.Client
}

func (provider *GitHubProvider) Initialize() {
	// transport := github.BasicAuthTransport{
	// 	Username: "x",
	// 	Password: "x",
	// }

	// provider.client = github.NewClient(transport.Client()) // TODO
	provider.client = github.NewClient(nil)

	RegisterProvider(provider)
}

func (provider *GitHubProvider) GetName() string {
	return "github"
}

func (provider *GitHubProvider) FetchReleases(owner string, repo string) ([]model.Release, error) {
	var releases []model.Release

	// TODO context.Background()
	ghReleases, _, err := provider.client.Repositories.ListReleases(context.Background(), owner, repo, nil)
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

func (provider *GitHubProvider) FetchTags(owner string, repo string) ([]model.Tag, error) {
	var tags []model.Tag

	// TODO context.Background()
	ghTags, _, err := provider.client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}

	for _, tag := range ghTags {
		url := tag.GetCommit().GetHTMLURL()
		if url == "" {
			url = fmt.Sprintf("https://github.com/%s/%s/commit/%s", owner, repo, tag.GetCommit().GetSHA())
		}

		tags = append(tags, model.Tag{
			Name:    tag.GetName(),
			URL:     url,
			Message: tag.GetCommit().GetMessage(),
		})
	}

	return tags, nil
}
