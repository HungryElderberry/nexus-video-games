package route

import (
	"nexus-video-games/internal/delivery/http/controller"
	"nexus-video-games/internal/delivery/http/middleware"
	"nexus-video-games/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteConfig struct {
	App                 *gin.Engine
	UserController      *controller.UserController
	BalanceController   *controller.BalanceController
	GameController      *controller.GameController
	SessionController   *controller.SessionController
	DashboardController *controller.DashboardController
	TokenService        *jwt.TokenService
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoutes()
	c.SetupAuthRoutes()
}

func (c *RouteConfig) SetupGuestRoutes() {
	// Interactive Swagger UI documentation
	c.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Game streaming landing route
	c.App.GET("/games/:slug", c.SessionController.GameLandingRoute)

	v1 := c.App.Group("/api/v1")
	{
		// Authentication & verification
		auth := v1.Group("/auth")
		{
			auth.POST("/register", c.UserController.Register)
			auth.GET("/verify", c.UserController.VerifyEmail)
			auth.POST("/login", c.UserController.Login)
			auth.POST("/reset-password", c.UserController.ResetPassword)
		}

		// Webhook callback (Xendit)
		v1.POST("/webhooks/xendit", c.BalanceController.HandleXenditWebhook)

		// Public game catalog exploration
		v1.GET("/games/explore", c.GameController.ExploreGames)
	}
}

func (c *RouteConfig) SetupAuthRoutes() {
	authMiddleware := middleware.AuthMiddleware(c.TokenService)
	v1 := c.App.Group("/api/v1")
	v1.Use(authMiddleware)
	{
		// User profile & Sign out
		v1.GET("/users/me", c.UserController.GetProfile)
		v1.POST("/auth/logout", c.UserController.Logout)

		// Wallet balance & topup
		balances := v1.Group("/balances")
		{
			balances.GET("", c.BalanceController.GetBalance)
			balances.POST("/topup", c.BalanceController.CreateTopup)
		}

		// Game management & library
		games := v1.Group("/games")
		{
			games.POST("/sync", c.GameController.SyncCatalog)
			games.POST("/buy", c.GameController.BuyGame)
			games.GET("/library", c.GameController.GetUserLibrary)
		}

		// Session management (Play 1 game at a time)
		sessions := v1.Group("/sessions")
		{
			sessions.POST("/launch/:game_id", c.SessionController.LaunchGame)
			sessions.GET("/active", c.SessionController.GetActiveSession)
		}

		// Player dashboard statistics
		v1.GET("/dashboard", c.DashboardController.GetDashboardStats)
	}
}
