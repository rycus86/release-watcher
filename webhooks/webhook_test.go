package webhooks

import (
	"encoding/json"
	"github.com/rycus86/release-watcher/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpWebhookSender(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		Project:  &mockProject{Name: "test.project"},
		Name:     "Example 1",
		Date:     time.Now(),
		URL:      "https://example.local/release/example/1",
		Webhooks: []string{
			server.URL + "/webhooks/one",
			server.URL + "/webhooks/two",
		},
	})

	sender.Send(&model.Release{
		Provider: mockProvider{Name: "test.provider"},
		Project:  &mockProject{Name: "test.project"},
		Name:     "Example 2",
		Date:     time.Now(),
		URL:      "https://example.local/release/example/2",
		Webhooks: []string{
			server.URL + "/webhooks/three",
		},
	})
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
	Name string
}

func (p *mockProject) String() string {
	return p.Name
}

func (p *mockProject) GetFilter() string {
	return ".*"
}
