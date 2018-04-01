package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
	"net/http"
)

type SlackNotificationManager struct {
	webhookUrl string
	channel    string
	username   string
	iconUrl    string

	httpClient *http.Client
}

func (m *SlackNotificationManager) initialize() {
	m.webhookUrl = env.Lookup("SLACK_WEBHOOK_URL", "/var/secrets/slack", "")
	m.channel = env.Lookup("SLACK_CHANNEL", "/var/secrets/slack", "")
	m.username = env.Lookup("SLACK_USERNAME", "/var/secrets/slack", "release-watcher")
	m.iconUrl = env.Lookup("SLACK_ICON_URL", "/var/secrets/slack", "")

	m.httpClient = &http.Client{
		Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/slack"),
		Transport: &transport.HttpTransportWithUserAgent{},
	}
}

func (m *SlackNotificationManager) Close() {
}

func (m *SlackNotificationManager) SendNotification(release *model.Release) error {
	if m.httpClient == nil {
		return errors.New("HTTP client not configured")
	}

	if m.webhookUrl == "" {
		return errors.New("webhook URL is not configured")
	}

	text := fmt.Sprintf("`[New release]` *%s* : %s", release.Project, release.Name)

	if release.URL != "" {
		text = fmt.Sprintf("%s\n%s", text, release.URL)
	}

	payload := map[string]string{
		"text": text,
	}

	if m.channel != "" {
		payload["channel"] = m.channel
	}

	if m.username != "" {
		payload["username"] = m.username
	}

	if m.iconUrl != "" {
		payload["icon_url"] = m.iconUrl
	}

	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response, err := m.httpClient.Post(m.webhookUrl, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("failed to send Slack message: %s", response.Status))
	}

	return nil
}
