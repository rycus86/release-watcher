package webhooks

import "github.com/rycus86/release-watcher/model"

type WebhookSender interface {
	Send(release *model.Release)
}

func NewWebhookSender() WebhookSender {
	return NewHttpWebhookSender()
}
