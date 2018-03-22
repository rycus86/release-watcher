package main

import (
	"fmt"
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/providers"
	"github.com/rycus86/release-watcher/store"
	"github.com/rycus86/release-watcher/watcher"
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
			log.Panic("Provider not found:", providerName)
		}

		for _, project := range projects {
			go WatchReleases(provider, project)
		}
	}
}

func WatchReleases(provider model.Provider, project model.Project) {
	rw, ok := provider.(watcher.ReleaseWatcher)
	if !ok {
		log.Println("The", provider.GetName(), "provider cannot watch releases")
		return
	}

	watcher.WatchReleases(rw, project, releaseChannel, shutdownChannel)
}

func WaitForChanges(db model.Store) {
	for {
		select {
		case releases := <-releaseChannel:
			lastKnown := ""
			hasNewRelease := false

			for _, release := range releases {
				if lastKnown == "" {  // TODO maybe keep all known releases in the database instead
					lastKnown = GetLastKnownRelease(release.Project, release.Provider, db)
				}

				if release.Name == lastKnown {
					break
				}

				// TODO proper filtering
				matched, err := regexp.MatchString("^[0-9]+\\.[0-9]+\\.[0-9]+$", release.Name)
				if !matched || err != nil {
					continue
				}

				log.Println(
					"[", release.Provider.GetName(), "]",
					"New release :", release.Project, ":", release.Name)

				if !hasNewRelease {
					hasNewRelease = true

					OnNewReleaseFound(release, db)
				}
			}

		case s := <-signalChannel:
			if s == syscall.SIGHUP {
				// TODO handle SIGHUP
			} else {
				close(shutdownChannel)
				return
			}
		}
	}
}

func GetLastKnownRelease(project model.Project, provider model.Provider, db model.Store) string {
	return db.Get(fmt.Sprintf("%s:%s:release", provider.GetName(), project))
}

func OnNewReleaseFound(release model.Release, db model.Store) {
	log.Println("Saving version", release.Name)
	db.Set(fmt.Sprintf("%s:%s:release", release.Provider.GetName(), release.Project), release.Name)

	// TODO notifications  https://github.com/nlopes/slack
}

func main() {
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbPath := config.Lookup("DATABASE_PATH", "/var/secrets/release-watcher", ":memory:")
	db, err := store.Initialize(dbPath)
	if err != nil {
		log.Panicln("Failed to initialize the database:", err)
	}
	defer db.Close()

	configPath := config.Lookup("CONFIGURATION_FILE", "/var/secrets/release-watcher", "release-watcher.yml")
	configuration, err := config.ParseConfigurationFile(configPath)
	if err != nil {
		log.Panicln("Failed load the configuration file:", err)
	}

	providers.InitializeProviders()

	StartWatchers(configuration)

	log.Println(
		"Started watching releases using",
		len(providers.GetProviders()), "providers",
	)

	WaitForChanges(db)

	log.Println("Application exiting")
}
