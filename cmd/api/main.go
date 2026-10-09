package main

import (
	"fmt"
	"log"

	"nexus-video-games/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	database := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	engine := config.NewGin(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       database,
		App:      engine,
		Validate: validate,
		Config:   viperConfig,
	})

	port := viperConfig.GetString("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Server] Starting Nexus Video Games API on port %s", port)
	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		panic(fmt.Sprintf("[Server] Engine failed to bind port: %v", err))
	}
}
