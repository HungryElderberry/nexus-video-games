package controller

import (
	"net/http"

	"nexus-video-games/internal/model"
	"nexus-video-games/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

// Register godoc
// @Summary      Register a new user account
// @Description  Creates a new user record and dispatches an email verification link via Resend.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterUserRequest true "User Registration Payload"
// @Success      201 {object} model.WebResponse[model.UserResponse]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/v1/auth/register [post]
func (ctrl *UserController) Register(c *gin.Context) {
	var req model.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := ctrl.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.WebResponse[*model.UserResponse]{
		Data: res,
	})
}

// VerifyEmail godoc
// @Summary      Verify user email address
// @Description  Validates the email token clicked from the Resend email message.
// @Tags         Authentication
// @Produce      json
// @Param        token query string true "Email verification token"
// @Success      200 {object} model.WebResponse[string]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/v1/auth/verify [get]
func (ctrl *UserController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: "Verification token query parameter is missing",
		})
		return
	}

	if err := ctrl.userUsecase.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Account verified successfully. You can now log in.",
	})
}

// Login godoc
// @Summary      Authenticate user and generate JWT
// @Description  Logs into account, returns Bearer token, and triggers security notification email.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body model.LoginUserRequest true "User Login Credentials"
// @Success      200 {object} model.WebResponse[model.LoginResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Router       /api/v1/auth/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var req model.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := ctrl.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.LoginResponse]{
		Data: res,
	})
}

// ResetPassword godoc
// @Summary      Reset forgotten or compromised password
// @Description  Sets a new password using the reset token provided in the login notification email.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body model.ResetPasswordRequest true "Password Reset Payload"
// @Success      200 {object} model.WebResponse[string]
// @Failure      400 {object} model.WebResponse[any]
// @Router       /api/v1/auth/reset-password [post]
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := ctrl.userUsecase.ResetPassword(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Password reset successful. Please sign in with your new credentials.",
	})
}

// Logout godoc
// @Summary      Sign out of the system
// @Description  Invalidates client-side session.
// @Tags         Authentication
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[string]
// @Failure      401 {object} model.WebResponse[any]
// @Router       /api/v1/auth/logout [post]
func (ctrl *UserController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, model.WebResponse[string]{
		Data: "Logged out successfully",
	})
}

// GetProfile godoc
// @Summary      Get current authenticated user profile
// @Description  Fetches the profile details of the user identified by the Bearer token.
// @Tags         Authentication
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} model.WebResponse[model.UserResponse]
// @Failure      401 {object} model.WebResponse[any]
// @Failure      404 {object} model.WebResponse[any]
// @Router       /api/v1/users/me [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.WebResponse[any]{
			Errors: "Unauthorized access",
		})
		return
	}
	userID := userIDVal.(uuid.UUID)

	res, err := ctrl.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.WebResponse[any]{
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.WebResponse[*model.UserResponse]{
		Data: res,
	})
}
