package providers

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/mitchellh/mapstructure"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
)

type GitHubProvider struct {
	client *github.Client
}

type GitHubProject struct {
	model.BaseProject `mapstructure:",squash"`

	Owner   string
	Repo    string
	UseTags bool
}

func (p GitHubProject) String() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Repo)
}

func (provider *GitHubProvider) Initialize() {
	token := env.Lookup("GITHUB_TOKEN", "/var/secrets/github", "")
	username := env.Lookup("GITHUB_USERNAME", "/var/secrets/github", "")
	password := env.Lookup("GITHUB_PASSWORD", "/var/secrets/github", "")

	if token != "" {
		ctx := context.Background()
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		oauthClient := oauth2.NewClient(ctx, tokenSource)

		provider.client = github.NewClient(oauthClient)

	} else if username != "" && password != "" {
		authenticatedTransport := github.BasicAuthTransport{
			Username: username,
			Password: password,

			Transport: &transport.HttpTransportWithUserAgent{},
		}

		provider.client = github.NewClient(authenticatedTransport.Client())

	} else {
		provider.client = github.NewClient(&http.Client{
			Transport: &transport.HttpTransportWithUserAgent{},
		})

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

	if project.UseTags {
		return provider.FetchTags(p)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/github"),
	)
	defer cancel()

	ghReleases, _, err := provider.client.Repositories.ListReleases(ctx, project.Owner, project.Repo, &github.ListOptions{PerPage: 50})
	if err != nil {
		return nil, err
	}

	for _, release := range ghReleases {
		name := release.GetName()
		if name == "" {
			name = release.GetTagName()
		}

		releases = append(releases, model.Release{
			Name: name,
			URL:  release.GetHTMLURL(),
			Date: release.GetPublishedAt().Time,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}

func (provider *GitHubProvider) FetchTags(p model.GenericProject) ([]model.Release, error) {
	var releases []model.Release

	project := p.(*GitHubProject)

	ctx, cancel := context.WithTimeout(
		context.Background(), env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/github"),
	)
	defer cancel()

	ghTags, _, err := provider.client.Repositories.ListTags(ctx, project.Owner, project.Repo, &github.ListOptions{PerPage: 50})
	if err != nil {
		return nil, err
	}

	for _, tag := range ghTags {
		name := tag.GetName()
		if name == "" {
			break
		}
		url := fmt.Sprintf("https://github.com/%v/%v/releases/tag/%v", project.Owner, project.Repo, name)
		releases = append(releases, model.Release{
			Name: name,
			URL:  url,

			Provider: provider,
			Project:  project,
		})
	}

	return releases, nil
}
