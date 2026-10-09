package config

import (
	"errors"
	"net/http"

	"nexus-video-games/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewGin(viper *viper.Viper) *gin.Engine {
	if viper.GetString("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(NewErrorHandler())

	return router
}

func NewErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if any errors occurred during the request lifecycle
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle validation errors specifically
			var valErr validator.ValidationErrors
			if errors.As(err, &valErr) {
				c.JSON(http.StatusBadRequest, model.WebResponse[any]{
					Errors: valErr.Error(),
				})
				return
			}

			// Fallback to internal server error if status code wasn't changed
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusInternalServerError
			}

			c.JSON(status, model.WebResponse[any]{
				Errors: err.Error(),
			})
		}
	}
}
