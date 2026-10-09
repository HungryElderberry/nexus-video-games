package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {

	_ = godotenv.Load()

	config := viper.New()
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		log.Printf("[Viper] Warning: .env file not found or failed to read: %v. Relying on environment variables.\n", err)
	}

	return config
}
