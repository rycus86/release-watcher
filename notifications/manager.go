package notifications

import "github.com/rycus86/release-watcher/env"
import "github.com/rycus86/release-watcher/model"

type NotificationManager interface {
	SendNotification(release *model.Release) error
	Close()
}

func NewNotificationManager() NotificationManager {
	service := env.Lookup("NOTIFICATION_SERVICE", "/var/secrets/slack", "slack")
	if service == "telegram" {
		manager := TelegramNotificationManager{}
		manager.initialize()
		return &manager
	} else {
		manager := SlackNotificationManager{}
		manager.initialize()
		return &manager
	}
}
