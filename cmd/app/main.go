package main

import (
	"effectiveMobile/pkg/config"
	"effectiveMobile/pkg/di"
	"log"
)

// @title People API
// @version 1.0
// @description This is a sample server People server.
// @contact.name API Support
// @contact.email support@example.com
// @host 0.0.0.0:8001
// @BasePath /
func main() {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(cfg)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
