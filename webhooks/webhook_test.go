package webhooks

import (
	"encoding/json"
	"github.com/rycus86/release-watcher/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHttpWebhookSender(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount += 1

		payload := WebhookPayload{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Error("Failed to decode the request")
		}

		if payload.Name == "Example 1" {
			if r.URL.Path != "/webhooks/one" && r.URL.Path != "/webhooks/two" {
				t.Error("Unexpected request:", r.URL)
			}

			if payload.URL != "https://example.local/release/example/1" {
				t.Error("Unexpected URL:", payload.URL)
			}
		} else if payload.Name == "Example 2" {
			if r.URL.Path != "/webhooks/three" {
				t.Error("Unexpected request:", r.URL)
			}

			if payload.URL != "https://example.local/release/example/2" {
				t.Error("Unexpected URL:", payload.URL)
			}
		} else {
			t.Error("Unexpected test case name:", payload.Name)
		}

		if payload.Provider != "test.provider" {
			t.Error("Unexpected provider:", payload.Provider)
		}
		if payload.Project != "test.project" {
			t.Error("Unexpected project:", payload.Project)
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	sender := NewHttpWebhookSender()

	sender.Send(&model.Release{
		Provider: mockProvider{Name: "test.provider"},
		Project: &mockProject{
			Name: "test.project",
			Webhooks: []string{
				server.URL + "/webhooks/one",
				server.URL + "/webhooks/two",
			},
		},
		Name: "Example 1",
		Date: time.Now(),
		URL:  "https://example.local/release/example/1",
	})

	sender.Send(&model.Release{
		Provider: mockProvider{Name: "test.provider"},
		Project: &mockProject{
			Name: "test.project",
			Webhooks: []string{
				server.URL + "/webhooks/three",
			},
		},
		Name: "Example 2",
		Date: time.Now(),
		URL:  "https://example.local/release/example/2",
	})

	time.Sleep(20 * time.Millisecond) // hack: give the goroutines some time

	if callCount != 3 {
		t.Error("Unexpected number of calls:", callCount)
	}
}

func TestHttpWebhookSenderWithAuthToken(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount += 1

		payload := WebhookPayload{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Error("Failed to decode the request")
		}

		if r.Header.Get("Webhook-Auth-Token") != "TestSecret" {
			t.Error("Invalid or missing auth token in", r.Header)
		}

		if r.URL.Path != "/webhooks/one" {
			t.Error("Unexpected request:", r.URL)
		}
		if payload.Name != "Example 1" {
			t.Error("Unexpected name:", payload.Name)
		}
		if payload.URL != "https://example.local/release/example/1" {
			t.Error("Unexpected URL:", payload.URL)
		}
		if payload.Provider != "test.provider" {
			t.Error("Unexpected provider:", payload.Provider)
		}
		if payload.Project != "test.project" {
			t.Error("Unexpected project:", payload.Project)
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	os.Setenv("HTTP_AUTHORIZATION", "TestSecret")
	defer os.Unsetenv("HTTP_AUTHORIZATION")

	sender := NewHttpWebhookSender()

	sender.Send(&model.Release{
		Provider: mockProvider{Name: "test.provider"},
		Project: &mockProject{
			Name: "test.project",
			Webhooks: []string{
				server.URL + "/webhooks/one",
			},
		},
		Name: "Example 1",
		Date: time.Now(),
		URL:  "https://example.local/release/example/1",
	})

	time.Sleep(20 * time.Millisecond) // hack: give the goroutines some time

	if callCount != 1 {
		t.Error("Unexpected number of calls:", callCount)
	}
}

type mockProvider struct {
	Name string
}

func (p mockProvider) Initialize() {

}

func (p mockProvider) GetName() string {
	return p.Name
}

func (p mockProvider) Parse(interface{}) model.GenericProject {
	return nil
}

type mockProject struct {
	Name     string
	Webhooks []string
}

func (p *mockProject) String() string {
	return p.Name
}

func (p *mockProject) GetFilter() string {
	return ".*"
}

func (p *mockProject) GetWebhooks() []string {
	return p.Webhooks
}
