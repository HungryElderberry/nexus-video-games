package usecase

import (
	"context"
	"errors"

	"nexus-video-games/internal/entity"
	"nexus-video-games/internal/gateway/email"
	"nexus-video-games/internal/model"
	"nexus-video-games/internal/pkg/jwt"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo     *postgresql.UserRepository
	balanceRepo  *postgresql.BalanceRepository
	tokenService *jwt.TokenService
	emailGateway *email.ResendGateway
}

func NewUserUsecase(
	userRepo *postgresql.UserRepository,
	balanceRepo *postgresql.BalanceRepository,
	tokenService *jwt.TokenService,
	emailGateway *email.ResendGateway,
) *UserUsecase {
	return &UserUsecase{
		userRepo:     userRepo,
		balanceRepo:  balanceRepo,
		tokenService: tokenService,
		emailGateway: emailGateway,
	}
}

func (u *UserUsecase) Register(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	existing, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Initialize empty wallet balance
	_ = u.balanceRepo.Create(ctx, &entity.UserBalance{
		UserID:        user.ID,
		BalanceAmount: 0,
		Currency:      "IDR",
	})

	// Dispatch email verification link
	verificationToken, err := u.tokenService.GenerateVerificationToken(user.ID, user.Email)
	if err == nil {
		_ = u.emailGateway.SendVerificationEmail(ctx, user.Email, verificationToken)
	}

	return &model.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}, nil
}

func (u *UserUsecase) VerifyEmail(ctx context.Context, token string) error {
	claims, err := u.tokenService.ValidateToken(token, jwt.TokenTypeVerification)
	if err != nil {
		return err
	}

	return u.userRepo.UpdateVerification(ctx, claims.UserID, true)
}

func (u *UserUsecase) Login(ctx context.Context, req *model.LoginUserRequest) (*model.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return nil, errors.New("account is not verified. Please check your email inbox")
	}

	accessToken, err := u.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Dispatch security notification email with emergency password-reset button
	resetToken, err := u.tokenService.GeneratePasswordResetToken(user.ID, user.Email)
	if err == nil {
		_ = u.emailGateway.SendLoginNotificationEmail(ctx, user.Email, resetToken)
	}

	return &model.LoginResponse{
		Token: accessToken,
		User: model.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

func (u *UserUsecase) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	claims, err := u.tokenService.ValidateToken(req.Token, jwt.TokenTypeReset)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(ctx, claims.UserID, string(hashedPassword))
}

func (u *UserUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}, nil
}
