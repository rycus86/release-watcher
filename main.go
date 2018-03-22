package main

import (
	"github.com/rycus86/release-watcher/config"
	"github.com/rycus86/release-watcher/model"
	"github.com/rycus86/release-watcher/providers"
	"github.com/rycus86/release-watcher/watcher"
	"github.com/rycus86/release-watcher/store"
	"log"
	"time"
	"os"
	"os/signal"
	"syscall"
)

var (
	releaseChannel = make(chan []model.Release)
	signalChannel = make(chan os.Signal, 1)
)

func Run(db model.Store, configuration *config.Configuration) {
	for providerName, projects := range configuration.Releases {
		provider := providers.GetProvider(providerName)
		if provider == nil {
			log.Panic("Provider not found:", providerName)
		}

		for _, project := range projects {
			go Watch(db, provider, project)
		}
	}

	// TODO tags (maybe with the same model as releases?)

	for {
		select {
		case releases := <-releaseChannel:
			for _, release := range releases {
				// TODO filtering
				Check(release)
			}

		case s := <-signalChannel:
			if s != syscall.SIGHUP { // TODO handle SIGHUP
				return
			}
		}
	}
}

func Watch(db model.Store, provider providers.Provider, project config.Project) {
	watcher, ok := provider.(watcher.ReleaseWatcher)
	if !ok {
		log.Println("The", provider.GetName(), "provider cannot watch releases")
		return
	}

	// TODO need to supply a different default here
	ticker := time.NewTicker(config.GetDuration("CHECK_INTERVAL", "/var/secrets/release-watcher"))

	for {
		select {
		case <- ticker.C:
			releases, err := watcher.FetchReleases(project)
			if err != nil {
				log.Println("Failed to fetch the releases of", project)
				continue
			}

			// TODO sort the releases here

			releaseChannel <- releases
		}

		// TODO shutdown/break
	}
}

func Check(release model.Release) {
	// TODO check
	// TODO will need the provider name too
}

func main() {
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbPath := config.Lookup("DATABASE_PATH", "/var/secrets/release-watcher", ":memory:")
	db, err := store.Initialize(dbPath)
	if err != nil {
		// TODO also include the error
		panic("Failed to initialize the database")
	}
	defer db.Close()

	configPath := config.Lookup("CONFIGURATION_FILE", "/var/secrets/release-watcher", "release-watcher.yml")
	configuration, err := config.ParseConfig(configPath)
	if err != nil {
		// TODO also include the error
		panic("Failed load the configuration file")
	}

	providers.InitializeProviders()

	log.Println(
		"Started watching", len(configuration.Releases), "projects for releases",
		"and", len(configuration.Tags), "projects for tags",
		"using", len(providers.GetProviders()), "providers",
	)

	Run(db, configuration)

	log.Println("Application exiting")
}
