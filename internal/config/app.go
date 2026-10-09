package config

import (
	"nexus-video-games/internal/delivery/http/controller"
	"nexus-video-games/internal/delivery/http/route"
	"nexus-video-games/internal/gateway/email"
	"nexus-video-games/internal/gateway/game"
	"nexus-video-games/internal/gateway/payment"
	"nexus-video-games/internal/pkg/jwt"
	"nexus-video-games/internal/repository/postgresql"
	"nexus-video-games/internal/usecase"

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
	// 1. Repositories
	userRepo := postgresql.NewUserRepository(config.DB)
	balanceRepo := postgresql.NewBalanceRepository(config.DB)
	topupRepo := postgresql.NewTopupRepository(config.DB)
	gameRepo := postgresql.NewGameRepository(config.DB)
	sessionRepo := postgresql.NewSessionRepository(config.DB)

	// 2. Security & Gateways
	tokenService := jwt.NewTokenService(config.Config)
	resendGateway := email.NewResendGateway(config.Config)
	xenditGateway := payment.NewXenditGateway(config.Config)
	itadGateway := game.NewITADGateway(config.Config)

	// 3. Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, balanceRepo, tokenService, resendGateway)
	balanceUsecase := usecase.NewBalanceUsecase(config.DB, balanceRepo, topupRepo, userRepo, xenditGateway, resendGateway)
	gameUsecase := usecase.NewGameUsecase(config.DB, gameRepo, balanceRepo, itadGateway)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, gameRepo, config.Config)
	dashboardUsecase := usecase.NewDashboardUsecase(balanceRepo, topupRepo, gameRepo)

	// 4. Controllers
	userController := controller.NewUserController(userUsecase)
	balanceController := controller.NewBalanceController(balanceUsecase)
	gameController := controller.NewGameController(gameUsecase)
	sessionController := controller.NewSessionController(sessionUsecase)
	dashboardController := controller.NewDashboardController(dashboardUsecase)

	// 5. Wire Routes
	routeConfig := route.RouteConfig{
		App:                 config.App,
		UserController:      userController,
		BalanceController:   balanceController,
		GameController:      gameController,
		SessionController:   sessionController,
		DashboardController: dashboardController,
		TokenService:        tokenService,
	}
	routeConfig.Setup()
}
