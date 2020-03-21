# Release Watcher

Watches new releases and sends notification when they're available.

Currently supports fetching releases from:

- [GitHub](https://github.com)
- [Docker Hub](https://hub.docker.com)
- [PyPI](TODO)
- [JetBrains](TODO)

Currently supports notifications to:

- [Slack](TODO)

The application periodically checks the target projects using their provider, and sends out a notificiation with some details about new releases when it finds them. Previously seen releases are kept in a [SQLite](TODO) database.

## Usage

With Docker:

```shell
$ docker run -d --name release-watcher 				\
	-v $PWD/release-watcher.yml:/var/conf/releases.yml:ro 	\
	-v $PWD/data:/data 					\
	-e CONFIGURATION_FILE=/var/conf/releases.yml 		\
	-e DATABASE_PATH=/data/releases.db			\
	-e SLACK_WEBHOOK_URL=https://some.slack.webhook.url	\
	rycus86/release-watcher
```

Alternatively, build with Go (tested on version 1.10), then execute:

```shell
$ export DATABASE_PATH=$PWD/releases.db
$ export SLACK_WEBHOOK_URL=https://some.slack.webhook.url
$ ./release-watcher
```

## Configuration

The application needs the target projects to be defined in a *YAML* configuration file. An example would be:

```yaml
releases:
  github:
    - owner: docker
      repo: docker-py

  dockerhub:
    - repo: nginx
    - owner: rycus86
      repo: grafana

  pypi:
    - name: flask

  jetbrains:
    - name: go
      alias: GoLand
```

The root `releases` holds mappings, keyed by the provider name, and the value being a list of project configurations. The available providers and their related configuration is listed below.

Each project can also have a `filter` property, which is a regular expression for whitelisting which releases to notify for. This defaults to `^[0-9]+.[0-9]+.[0-9]+$`, so it would match `1.2.3`, but not `1.2.3-rc1` for example.

### GitHub

The `github` provider looks for GitHub releases. The configuration items need to have an `owner` and `repo` property.

This provider accepts some optional configuration parameters, either from environment variables, or a key-value file at `/var/secrets/github` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `GITHUB_TOKEN` | The GitHub OAuth token if you wish to make authenticated API calls, to get better rate-limiting. ([token generation](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) - no scopes are needed) | - |
| `GITHUB_USERNAME` | The GitHub username for basic authentication instead of Oauth with token | - |
| `GITHUB_PASSWORD` | The GitHub password for the username | - |
| `HTTP_TIMEOUT` | The HTTP timeout for API calls | `30s` |

### Docker Hub

The `dockerhub` provider looks for new image tags on Docker Hub. The configuration items need to have a `repo` property and usually an `owner`, which defaults to `library` if omitted.

This provider accepts some optional configuration parameters, either from environment variables, or a key-value file at `/var/secrets/dockerhub` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `HTTP_TIMEOUT` | The HTTP timeout for API calls | `30s` |
| `PAGE_SIZE` | The number of tags to check on each call, per project | `50` |

### PyPI

The `pypi` provider looks for new releases on PyPI. The configuration items need to have a `name` property.

This provider accepts some optional configuration parameters, either from environment variables, or a key-value file at `/var/secrets/pypi` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `HTTP_TIMEOUT` | The HTTP timeout for API calls | `30s` |

### JetBrains

The `jetbrains` provider looks for new releases for JetBrains applications. The configuration items need to have a `name` property, which is the application ID as JetBrains knows it. It could also include an `alias` property to display in the notifications instead of the ID.

This provider accepts some optional configuration parameters, either from environment variables, or a key-value file at `/var/secrets/jetbrains` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `HTTP_TIMEOUT` | The HTTP timeout for API calls | `30s` |

### Slack notifications

The Slack notification manager accepts some optional configuration parameters, either from environment variables, or a key-value file at `/var/secrets/slack` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `SLACK_WEBHOOK_URL` | The target Slack webhook to send the notifications to __(required)__ | - |
| `SLACK_CHANNEL` | The target channel to use, otherwise the default for the webhook | - |
| `SLACK_USERNAME` | The user name to display on the notification in Slack | `release-watcher` |
| `SLACK_ICON_URL` | The URL of the image to display on the notification | - |
| `HTTP_TIMEOUT` | The HTTP timeout for API calls | `30s` |

### General configuration

There are some more, general configuration parameters, either from environment variables, or a key-value file at `/var/secrets/release-watcher` in `KEY=VALUE` format on each line.

| Key | Description | Default |
| --- | ----------- | ------- |
| `CONFIGURATION_FILE` | The configuration *YAML* file for the projects | `release-watcher.yml` |
| `DATABASE_PATH` | The path to the SQLite database file | `file::memory:?cache=shared` |
| `CHECK_INTERVAL` | The frequency to check for new releases | `4h` |

## License

MIT
