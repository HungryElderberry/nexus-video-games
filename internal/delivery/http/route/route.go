package route

import (
	"net/http"

	"nexus-video-games/internal/model"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App *gin.Engine
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoutes()
	c.SetupAuthRoutes()
}

func (c *RouteConfig) SetupGuestRoutes() {
	v1 := c.App.Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, model.WebResponse[string]{
				Data: "Nexus Video Games API engine operational",
			})
		})
	}
}

func (c *RouteConfig) SetupAuthRoutes() {

	// TODO

}
