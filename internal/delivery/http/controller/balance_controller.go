package controller

import (
	"net/http"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BalanceController struct {
	balanceUsecase *usecase.BalanceUsecase
}

func NewBalanceController(balanceUsecase *usecase.BalanceUsecase) *BalanceController {
	return &BalanceController{balanceUsecase: balanceUsecase}
}

// GetBalance godoc
// @Summary      Get current in-game wallet balance
// @Description  Returns wallet amount in minor units (IDR) for the authenticated player.
// @Tags         Balances & Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[model.BalanceResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/balances [get]
func (ctrl *BalanceController) GetBalance(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.WebResponse[any]{
			Errors: "Unauthorized access",
		})
		return
	}
	userID := userIDVal.(uuid.UUID)

	res, err := ctrl.balanceUsecase.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.BalanceResponse]{
		Data: res,
	})
}

// CreateTopup godoc
// @Summary      Create Xendit invoice for wallet top-up
// @Description  Generates a checkout URL where the player can simulate or complete payment.
// @Tags         Balances & Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body model.TopupBalanceRequest true "Top-up Amount"
// @Success      201 {object} model.WebResponse[model.TopupResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      500 {object} model.WebResponse[any]
// @Router       /api/v1/balances/topup [post]
func (ctrl *BalanceController) CreateTopup(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	emailVal, _ := c.Get("user_email")

	userID := userIDVal.(uuid.UUID)
	email := emailVal.(string)

	var req model.TopupBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := ctrl.balanceUsecase.CreateTopup(c.Request.Context(), userID, email, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.WebResponse[*model.TopupResponse]{
		Data: res,
	})
}

// HandleXenditWebhook godoc
// @Summary      Handle incoming payment webhooks from Xendit
// @Description  Processes payment settlement notifications, credits player wallet, and triggers email receipt.
// @Tags         Balances & Payments
// @Accept       json
// @Produce      json
// @Param        x-callback-token header string false "Xendit Webhook Verification Token"
// @Param        payload body model.XenditInvoiceWebhookPayload true "Xendit Webhook Payload"
// @Success      200 {object} model.WebResponse[string]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/v1/webhooks/xendit [post]
func (ctrl *BalanceController) HandleXenditWebhook(c *gin.Context) {
	webhookToken := c.GetHeader("x-callback-token")

	var payload model.XenditInvoiceWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}

	if err := ctrl.balanceUsecase.HandleXenditWebhook(c.Request.Context(), webhookToken, &payload); err != nil {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Webhook received and processed",
	})
}
