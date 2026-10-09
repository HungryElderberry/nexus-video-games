package controller

import (
	"fmt"
	"net/http"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionController struct {
	sessionUsecase *usecase.SessionUsecase
}

func NewSessionController(sessionUsecase *usecase.SessionUsecase) *SessionController {
	return &SessionController{sessionUsecase: sessionUsecase}
}

// LaunchGame godoc
// @Summary      Launch a game and start active playing session
// @Description  Verifies ownership, enforces single-game active session rule (with 60-second auto-expiry), and provides game redirect route.
// @Tags         Game Sessions
// @Security     BearerAuth
// @Produce      json
// @Param        game_id path string true "External Game ID"
// @Success      200 {object} model.WebResponse[model.LaunchGameResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      409 {object} model.WebResponse[any]
// @Router       /api/v1/sessions/launch/{game_id} [post]
func (ctrl *SessionController) LaunchGame(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	gameID := c.Param("game_id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: "Game ID parameter is required",
		})
		return
	}

	res, err := ctrl.sessionUsecase.LaunchGame(c.Request.Context(), userID, gameID)
	if err != nil {
		c.JSON(http.StatusConflict, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.LaunchGameResponse]{
		Data: res,
	})
}

// GetActiveSession godoc
// @Summary      Check current active playing session
// @Description  Returns details of the currently running session and countdown to expiry.
// @Tags         Game Sessions
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[model.ActiveSessionResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      404 {object} model.WebResponse[any]
// @Router       /api/v1/sessions/active [get]
func (ctrl *SessionController) GetActiveSession(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	session, err := ctrl.sessionUsecase.GetActiveSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.ActiveSessionResponse]{
		Data: session,
	})
}

// GameLandingRoute godoc
// @Summary      Simulated game streaming route
// @Description  Destination route for game launching.
// @Tags         Game Sessions
// @Produce      json
// @Param        slug path string true "Game Slug / External Game ID"
// @Success      200 {object} model.WebResponse[string]
// @Router       /games/{slug} [get]
func (ctrl *SessionController) GameLandingRoute(c *gin.Context) {
	slug := c.Param("slug")
	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: fmt.Sprintf("Nexus Game Launcher Engine: Now streaming and playing game '%s'", slug),
	})
}
