package controller

import (
	"net/http"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardController struct {
	dashboardUsecase *usecase.DashboardUsecase
}

func NewDashboardController(dashboardUsecase *usecase.DashboardUsecase) *DashboardController {
	return &DashboardController{dashboardUsecase: dashboardUsecase}
}

// GetDashboardStats godoc
// @Summary      Get player dashboard metrics
// @Description  Calculates user wallet balance, owned games count, lifetime spending, monthly spending, and games played in the current month.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[model.DashboardStatsResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/dashboard [get]
func (ctrl *DashboardController) GetDashboardStats(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.WebResponse[any]{
			Errors: "Unauthorized access",
		})
		return
	}
	userID := userIDVal.(uuid.UUID)

	stats, err := ctrl.dashboardUsecase.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.DashboardStatsResponse]{
		Data: stats,
	})
}
