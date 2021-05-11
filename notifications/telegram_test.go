package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
)

func TestSendTelegramNotification(t *testing.T) {
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

		if r.UserAgent() != transport.DefaultUserAgent {
			t.Error("Unexpected User-Agent:", r.UserAgent())
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	os.Setenv("NOTIFICATION_SERVICE", "telegram")
	defer os.Unsetenv("NOTIFICATION_SERVICE")
	os.Setenv("TELEGRAM_BOT_TOKEN", `110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw`)
	defer os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Setenv("TELEGRAM_CHAT_ID", "1234567890")
	defer os.Unsetenv("TELEGRAM_CHAT_ID")

	manager := NewNotificationManager()
	defer manager.Close()

	release := model.Release{
		Project: &mockProject{},

		Name: "0.1.2",
	}

	err := manager.SendNotification(&release)
	if err != nil {
		t.Error("Failed to send notification:", err)
	}
}
