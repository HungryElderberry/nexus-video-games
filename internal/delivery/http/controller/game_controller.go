package controller

import (
	"net/http"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GameController struct {
	gameUsecase *usecase.GameUsecase
}

func NewGameController(gameUsecase *usecase.GameUsecase) *GameController {
	return &GameController{gameUsecase: gameUsecase}
}

// ExploreGames godoc
// @Summary      List all available games in catalog
// @Description  Retrieves games cached from IsThereAnyDeal available for wallet purchase.
// @Tags         Games
// @Produce      json
// @Success      200 {object} model.WebResponse[[]model.GameCatalogItemResponse]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/games/explore [get]
func (ctrl *GameController) ExploreGames(c *gin.Context) {
	games, err := ctrl.gameUsecase.ExploreGames(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[[]model.GameCatalogItemResponse]{
		Data: games,
	})
}

// SyncCatalog godoc
// @Summary      Synchronize catalog from IsThereAnyDeal API
// @Description  Pulls live game deals and upserts records into PostgreSQL games_catalog.
// @Tags         Games
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[string]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/games/sync [post]
func (ctrl *GameController) SyncCatalog(c *gin.Context) {
	if err := ctrl.gameUsecase.SyncCatalog(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Game catalog successfully updated from IsThereAnyDeal",
	})
}

// BuyGame godoc
// @Summary      Purchase a game using wallet balance
// @Description  Deducts game price from user wallet and grants game ownership.
// @Tags         Games
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.BuyGameRequest true "External Game ID to purchase"
// @Success      200 {object} model.WebResponse[string]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Router       /api/v1/games/buy [post]
func (ctrl *GameController) BuyGame(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	var req model.BuyGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := ctrl.gameUsecase.BuyGame(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Game purchase completed. Added to your personal library.",
	})
}

// GetUserLibrary godoc
// @Summary      Get games owned by authenticated user
// @Description  Lists all acquired games along with logged playtime.
// @Tags         Games
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[[]model.UserGameLibraryResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/games/library [get]
func (ctrl *GameController) GetUserLibrary(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	library, err := ctrl.gameUsecase.GetUserLibrary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[[]model.UserGameLibraryResponse]{
		Data: library,
	})
}
