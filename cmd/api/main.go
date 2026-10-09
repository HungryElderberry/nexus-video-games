package main

import (
	"fmt"
	"log"

	_ "nexus-video-games/docs"
	"nexus-video-games/internal/config"
)

// @title                      Nexus Video Games REST API
// @version                    1.0
// @description                High-performance backend API featuring video games catalog, in-game wallet, Xendit payment gateway, Resend email alerts, and single-game session enforcement.
// @host                       localhost:8080
// @BasePath                   /
// @schemes                    http https

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Enter your token with the format: Bearer <your_jwt_token>

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

	serverURL := fmt.Sprintf("http://localhost:%s", port)
	swaggerURL := fmt.Sprintf("http://localhost:%s/swagger/index.html", port)

	log.Printf("[Server] Starting Nexus Video Games API...")
	log.Printf("[Server] Base URL:  %s", serverURL)
	log.Printf("[Swagger] UI Docs:  %s", swaggerURL)

	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		panic(fmt.Sprintf("[Server] Engine failed to bind port: %v", err))
	}
}
