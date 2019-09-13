package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/transport"
	"io"
	"net/http"
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
	if release.Webhooks == nil {
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

	for _, webhookTargetUrl := range release.Webhooks {
		go sendMessage(webhookTargetUrl, s.client, body)
	}
}

func sendMessage(targetUrl string, client *http.Client, jsonBody io.Reader) {
	if response, err := client.Post(targetUrl, "application/json", jsonBody); err != nil {
		fmt.Println("Failed to POST a webhook to", targetUrl, ":", err)
	} else {
		defer response.Body.Close()

		if response.StatusCode >= 400 {
			fmt.Println("Failed to POST a webhook to", targetUrl, ":", err)
		} else {
			fmt.Println("Webhook successfully POSTed to", targetUrl)
		}
	}
}
