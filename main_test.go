package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/providers"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestStartWatchers(t *testing.T) {
	shutdownChannel = make(chan struct{})
	defer close(shutdownChannel)

	registerResponderFromFile(
		"https://api.github.com/repos/docker/docker-py/releases",
		"./testdata/github_releases.json",
	)

	configuration := model.Configuration{
		Releases: map[string][]model.GenericProject{
			"github": {
				&providers.GitHubProject{Owner: "docker", Repo: "docker-py"},
			},
		},
	}

	StartWatchers(&configuration)

	releases := <-releaseChannel
	if len(releases) != 30 {
		t.Error("Expected to find 30 releases")
	}

	if releases[0].Name != "3.1.3" {
		t.Error("Unexpected release:", releases[0].Name)
	}
}

func TestPanicOnUnknownProviders(t *testing.T) {
	shutdownChannel = make(chan struct{})
	defer close(shutdownChannel)

	defer func() {
		err := recover()

		if err == nil {
			t.Error("Expected to panic")
		}
		if !strings.Contains(err.(string), "testfake") {
			t.Error("Unexpected error message:", err)
		}
	}()

	configuration := model.Configuration{
		Releases: map[string][]model.GenericProject{
			"github":   {},
			"testfake": {},
		},
	}

	StartWatchers(&configuration)
}

func TestWaitForChanges(t *testing.T) {
	shutdownChannel = make(chan struct{})

	registerResponderFromFile(
		"https://api.github.com/repos/docker/docker-py/releases",
		"./testdata/github_releases.json",
	)
	registerResponderFromFile(
		"https://hub.docker.com/v2/repositories/rycus86/grafana/tags/",
		"./testdata/dockerhub_releases.json",
	)
	registerResponderFromFile(
		"https://pypi.python.org/pypi/prometheus-flask-exporter/json",
		"./testdata/pypi_releases.json",
	)
	registerResponderFromFile(
		"https://data.services.jetbrains.com/products?code=GO",
		"./testdata/jetbrains_releases.json",
	)

	configuration := model.Configuration{
		Releases: map[string][]model.GenericProject{
			"github": {
				&providers.GitHubProject{Owner: "docker", Repo: "docker-py"},
			},
			"dockerhub": {
				&providers.DockerHubProject{Owner: "rycus86", Repo: "grafana"},
			},
			"pypi": {
				&providers.PyPIProject{Name: "prometheus-flask-exporter"},
			},
			"jetbrains": {
				&providers.JetBrainsProject{Name: "go", Alias: "GoLand"},
			},
		},
	}

	StartWatchers(&configuration)

	go func() {
		time.Sleep(20 * time.Millisecond)
		signalChannel <- syscall.SIGTERM
	}()

	store := mockStore{}
	notifier := mockNotifier{}
	webhookSender := mockWebhookSender{}

	WaitForChanges(&store, &notifier, &webhookSender, nil)

	if !strings.Contains(store.callExists, "docker/docker-py:3.1.3") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callExists, "docker/docker-py:2.7.0") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callMark, "docker/docker-py:3.1.3") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}
	if !strings.Contains(store.callMark, "docker/docker-py:2.7.0") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}

	if strings.Contains(notifier.callSend, "docker/docker-py:2.7.0") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}
	if !strings.Contains(notifier.callSend, "docker/docker-py:3.1.3") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}

	if !strings.Contains(webhookSender.callSend, "docker/docker-py:2.7.0") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}
	if !strings.Contains(webhookSender.callSend, "docker/docker-py:3.1.3") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}

	if strings.Contains(store.callExists, "rycus86/grafana:latest") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callExists, "rycus86/grafana:5.0.2") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callExists, "rycus86/grafana:4.6.1") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callMark, "rycus86/grafana:5.0.2") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}
	if !strings.Contains(store.callMark, "rycus86/grafana:4.6.1") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}

	if strings.Contains(notifier.callSend, "rycus86/grafana:4.6.1") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}
	if !strings.Contains(notifier.callSend, "rycus86/grafana:5.0.2") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}

	if !strings.Contains(webhookSender.callSend, "rycus86/grafana:4.6.1") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}
	if !strings.Contains(webhookSender.callSend, "rycus86/grafana:5.0.2") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}

	if !strings.Contains(store.callExists, "prometheus-flask-exporter:0.2.1") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callExists, "prometheus-flask-exporter:0.1.0") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callMark, "prometheus-flask-exporter:0.2.1") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}
	if !strings.Contains(store.callMark, "prometheus-flask-exporter:0.1.0") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}

	if strings.Contains(notifier.callSend, "prometheus-flask-exporter:0.1.0") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}
	if !strings.Contains(notifier.callSend, "prometheus-flask-exporter:0.2.1") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}

	if !strings.Contains(webhookSender.callSend, "prometheus-flask-exporter:0.1.0") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}
	if !strings.Contains(webhookSender.callSend, "prometheus-flask-exporter:0.2.1") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}

	if !strings.Contains(store.callExists, "GoLand:2018.1 (181.4203.567)") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callExists, "GoLand:2018.1 (181.4203.544 eap)") {
		t.Error("Unexpected calls to store.Exists:", store.callExists)
	}
	if !strings.Contains(store.callMark, "GoLand:2018.1 (181.4203.567)") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}
	if !strings.Contains(store.callMark, "GoLand:2018.1 (181.4203.544 eap)") {
		t.Error("Unexpected calls to store.Mark:", store.callMark)
	}

	if strings.Contains(notifier.callSend, "GoLand:2018.1 (181.4203.544 eap)") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}
	if !strings.Contains(notifier.callSend, "GoLand:2018.1 (181.4203.567)") {
		t.Error("Unexpected calls to notifier.SendNotification:", notifier.callSend)
	}

	if !strings.Contains(webhookSender.callSend, "GoLand:2018.1 (181.4203.544 eap)") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}
	if !strings.Contains(webhookSender.callSend, "GoLand:2018.1 (181.4203.567)") {
		t.Error("Unexpected calls to webhookSender.Send:", notifier.callSend)
	}
}

func TestReload(t *testing.T) {
	shutdownChannel = make(chan struct{})

	go func() {
		time.Sleep(20 * time.Millisecond)
		signalChannel <- syscall.SIGHUP
		signalChannel <- syscall.SIGHUP
		signalChannel <- syscall.SIGTERM
	}()

	store := mockStore{}
	notifier := mockNotifier{}
	webhookSender := mockWebhookSender{}

	reloaded := 0
	reloader := func() {
		reloaded++
	}

	WaitForChanges(&store, &notifier, &webhookSender, reloader)

	if reloaded != 2 {
		t.Error("Unexpected number of reloads:", reloaded)
	}
}

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// TODO how to deal with the shutdownChannel here?

	providers.InitializeProviders()

	os.Exit(m.Run())
}

func registerResponderFromFile(url string, path string) {
	testdata, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Failed to load test data")
	}

	registerResponder(url, string(testdata))
}

func registerResponder(url string, response string) {
	httpmock.RegisterResponder(
		"GET", url,
		httpmock.NewStringResponder(200, response),
	)
}

type mockStore struct {
	callExists string
	callMark   string
}

func (s *mockStore) Exists(r model.Release) bool {
	s.callExists = s.callExists + "," + fmt.Sprintf("%s:%s", r.Project.String(), r.Name)
	return false
}

func (s *mockStore) Mark(r model.Release) error {
	s.callMark = s.callMark + "," + fmt.Sprintf("%s:%s", r.Project.String(), r.Name)
	return nil
}

func (s *mockStore) Close() {
}

type mockNotifier struct {
	callSend string
}

func (n *mockNotifier) SendNotification(r *model.Release) error {
	n.callSend = n.callSend + "," + fmt.Sprintf("%s:%s", r.Project.String(), r.Name)
	return nil
}

func (n *mockNotifier) Close() {
}

type mockWebhookSender struct {
	callSend string
}

func (s *mockWebhookSender) Send(r *model.Release) {
	s.callSend = s.callSend + "," + fmt.Sprintf("%s:%s", r.Project.String(), r.Name)
}
