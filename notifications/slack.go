package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
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
	m.webhookUrl = config.Lookup("SLACK_WEBHOOK_URL", "/var/secrets/slack", "")
	m.channel = config.Lookup("SLACK_CHANNEL", "/var/secrets/slack", "")
	m.username = config.Lookup("SLACK_USERNAME", "/var/secrets/slack", "release-watcher")
	m.iconUrl = config.Lookup("SLACK_ICON_URL", "/var/secrets/slack", "")

	m.httpClient = &http.Client{
		Timeout: config.GetTimeout("HTTP_TIMEOUT", "/var/secrets/slack"),
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

	payload := map[string]string{
		"text": fmt.Sprintf("`[New release]` *%s* : %s", release.Project, release.Name),
	}

	if m.channel != "" {
		payload["channel"] = m.channel
	}

	if m.username != "" {
		payload["username"] = m.channel
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
