package middleware

import (
	"net/http"
	"strings"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService *jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.WebResponse[any]{
				Errors: "Authorization header missing",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.WebResponse[any]{
				Errors: "Authorization format must be Bearer {token}",
			})
			return
		}

		claims, err := tokenService.ValidateToken(parts[1], jwt.TokenTypeAccess)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.WebResponse[any]{
				Errors: err.Error(),
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
