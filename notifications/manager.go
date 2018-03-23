package notifications

import "github.com/rycus86/release-watcher/model"

type NotificationManager interface {
	SendNotification(release *model.Release) error
	Close()
}

func NewNotificationManager() NotificationManager {
	manager := SlackNotificationManager{}
	manager.initialize()
	return &manager
}
