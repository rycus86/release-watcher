package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
)

type HttpWebhookSender struct {
	client *http.Client

	authorizationKey string
}

func NewHttpWebhookSender() *HttpWebhookSender {
	return &HttpWebhookSender{
		client: &http.Client{
			Timeout:   env.GetTimeout("HTTP_TIMEOUT", "/var/secrets/webhooks"),
			Transport: &transport.HttpTransportWithUserAgent{},
		},
		authorizationKey: env.Lookup("HTTP_AUTHORIZATION", "/var/secrets/webhooks", ""),
	}
}

func (s *HttpWebhookSender) Send(release *model.Release) {
	if release.Project.GetWebhooks() == nil {
		return
	}

	payload := WebhookPayload{
		Provider: release.Provider.GetName(),
		Project:  release.Project.String(),
		Name:     release.Name,
		Date:     release.Date,
		URL:      release.URL,
	}

	var body = new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		fmt.Printf("Invalid webhook payload: %+v -- %s\n", payload, err)
		return
	}

	for _, webhookTargetUrl := range release.Project.GetWebhooks() {
		go s.sendMessage(webhookTargetUrl, bytes.NewReader(body.Bytes()))
	}
}

func (s *HttpWebhookSender) sendMessage(targetUrl string, jsonBody io.Reader) {
	request, err := http.NewRequest("POST", targetUrl, jsonBody)
	if err != nil {
		fmt.Println("Failed to prepare a request to", targetUrl, ":", err)
		return
	}

	request.Header.Add("Content-Type", "application/json")

	if key := s.authorizationKey; key != "" {
		request.Header.Add("Webhook-Auth-Token", key)
	}

	if response, err := s.client.Do(request); err != nil {
		fmt.Println("Failed to POST a webhook to", targetUrl, ":", err)
	} else {
		defer response.Body.Close()

		if response.StatusCode >= 400 {
			fmt.Println("Failed to POST a webhook to", targetUrl, ":", response.Status)
		} else {
			fmt.Println("Webhook successfully POSTed to", targetUrl)
		}
	}
}
