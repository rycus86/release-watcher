package main

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/env"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/notifications"
	"github.com/rycus86/release-watcher/providers"
	"github.com/rycus86/release-watcher/store"
	"github.com/rycus86/release-watcher/watcher"
	"github.com/rycus86/release-watcher/webhooks"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var (
	releaseChannel  = make(chan []model.Release)
	signalChannel   = make(chan os.Signal, 1)
	shutdownChannel = make(chan struct{})
)

func StartWatchers(configuration *model.Configuration) {
	for providerName, projects := range configuration.Releases {
		provider := providers.GetProvider(providerName)

		if provider == nil {
			log.Panicln("Provider not found:", providerName)
		}

		for _, project := range projects {
			go WatchReleases(provider, project)
		}
	}
}

func WatchReleases(provider model.Provider, project model.GenericProject) {
	rw, ok := provider.(watcher.ReleaseWatcher)
	if !ok {
		log.Println("The", provider.GetName(), "provider cannot watch releases")
		return
	}

	watcher.WatchReleases(rw, project, releaseChannel, shutdownChannel)
}

func WaitForChanges(
	db model.Store,
	notifier notifications.NotificationManager,
	webhookSender webhooks.WebhookSender,
	reloadHandler func()) {

	for {
		select {
		case releases := <-releaseChannel:
			hasNewRelease := false

			for _, release := range releases {
				filter := release.Project.GetFilter()

				matched, err := regexp.MatchString(filter, release.Name)
				if !matched || err != nil {
					continue
				}

				if db.Exists(release) {
					break
				}

				log.Println(
					"[", release.Provider.GetName(), "]",
					"New release :", release.Project, ":", release.Name)

				if err := db.Mark(release); err != nil {
					log.Println("Failed to save the new version:", err)
				}

				webhookSender.Send(&release)

				if !hasNewRelease {
					hasNewRelease = true

					if err := notifier.SendNotification(&release); err != nil {
						log.Println("Failed to send notifications:", err)
					}
				}
			}

		case s := <-signalChannel:
			if s == syscall.SIGHUP {
				log.Println("Reload signal received")

				reloadHandler()

			} else {
				log.Println("Shutdown signal received")

				close(shutdownChannel)
				return

			}
		}
	}
}

func main() {
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbPath := env.Lookup("DATABASE_PATH", "/var/secrets/release-watcher", "file::memory:?cache=shared")
	db, err := store.Initialize(dbPath)
	if err != nil {
		log.Panicln("Failed to initialize the database:", err)
	}
	defer db.Close()

	providers.InitializeProviders()

	configPath := env.Lookup("CONFIGURATION_FILE", "/var/secrets/release-watcher", "release-watcher.yml")
	configuration, err := config.ParseConfigurationFile(configPath)
	if err != nil {
		log.Panicln("Failed load the configuration file:", err)
	}

	notifier := notifications.NewNotificationManager()
	defer notifier.Close()

	webhookSender := webhooks.NewWebhookSender()

	reloadHandler := func() {
		close(shutdownChannel)

		shutdownChannel = make(chan struct{})

		config.Reload(configuration)
		StartWatchers(configuration)
	}

	StartWatchers(configuration)

	log.Println(
		"Started watching releases using",
		len(providers.GetProviders()), "providers",
	)

	WaitForChanges(db, notifier, webhookSender, reloadHandler)

	log.Println("Application exiting")
}
