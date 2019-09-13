package webhooks

import "time"

type WebhookPayload struct {
	Provider string
	Project  string

	Name string
	Date time.Time
	URL  string
}
