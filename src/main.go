package main

import (
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"github.com/k0rean-rand0m/event-horizon/src/db"
	"github.com/k0rean-rand0m/event-horizon/src/utils"
	"github.com/k0rean-rand0m/event-horizon/src/watcher"
	"log"
)

func main() {
	// Preparing configuration
	config.Setup()

	// Initializing database
	if config.Config.DbUrl != "" {
		err := db.Instance.Setup()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Running watcher
	err := watcher.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Waiting for all goroutines to finish (due to reconnection mechanism is not expected)
	utils.Wg.Wait()
}
