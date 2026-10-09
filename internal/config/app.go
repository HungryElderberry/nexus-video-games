package config

import (
	"nexus-video-games/internal/delivery/http/route"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *gin.Engine
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// 1. Repositories initialization
	// 2. Gateways initialization (Xendit, Resend, ITAD)
	// 3. Usecases initialization
	// 4. Controllers initialization

	// 5. Route registration
	routeConfig := route.RouteConfig{
		App: config.App,
	}
	routeConfig.Setup()
}
