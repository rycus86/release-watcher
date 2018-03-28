package notifications

import (
	"github.com/rycus86/release-watcher/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Error("Failed to decode the request")
		}

		text, ok := payload["text"]
		if !ok {
			t.Error("Missing 'text' field")
		}
		if text != "`[New release]` *test/repo* : 0.1.2" {
			t.Error("Unexpected text:", text)
		}

		if channel, ok := payload["channel"]; !ok || channel != "#testing" {
			t.Error("Unexpected channel:", channel)
		}
		if username, ok := payload["username"]; !ok || username != "Test Runner" {
			t.Error("Unexpected username:", username)
		}
		if icon, ok := payload["icon_url"]; !ok || icon != "http://slack.icon" {
			t.Error("Unexpected icon URL:", icon)
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	os.Setenv("SLACK_WEBHOOK_URL", server.URL)
	defer os.Unsetenv("SLACK_WEBHOOK_URL")
	os.Setenv("SLACK_CHANNEL", "#testing")
	defer os.Unsetenv("SLACK_CHANNEL")
	os.Setenv("SLACK_USERNAME", "Test Runner")
	defer os.Unsetenv("SLACK_USERNAME")
	os.Setenv("SLACK_ICON_URL", "http://slack.icon")
	defer os.Unsetenv("SLACK_ICON_URL")

	manager := NewNotificationManager()
	defer manager.Close()

	release := model.Release{
		Project: model.Project{
			Owner: "test",
			Repo:  "repo",
		},

		Name: "0.1.2",
	}

	err := manager.SendNotification(&release)
	if err != nil {
		t.Error("Failed to send notification:", err)
	}
}

func TestDoesNotSendWithoutWebhookUrl(t *testing.T) {
	manager := NewNotificationManager()
	defer manager.Close()

	release := model.Release{
		Project: model.Project{
			Owner: "test",
			Repo:  "repo",
		},

		Name: "0.1.2",
	}

	err := manager.SendNotification(&release)
	if err == nil {
		t.Error("Expected to fail")
	}
}
