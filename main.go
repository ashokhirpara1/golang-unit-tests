package main

import (
	"unit-tests/configuration"
	"unit-tests/internal"
	"unit-tests/mongo"
	"unit-tests/prettylogs"
	"unit-tests/routes"
)

// main function to boot up everything
func main() {

	// Get application specific configurations
	config := configuration.Get()

	// Initialize structured logs
	logs := prettylogs.Get()

	db, err := mongo.Get(config, logs)
	if err != nil {
		logs.Fatal("database.Get", "Failed to initialize database", err)
		return
	}

	ctlr := internal.InitController(db, logs, config)

	logs.Info("main function to boot up everything")

	// Creates a http server
	routes.InitRoutes(ctlr, logs, config)
}
