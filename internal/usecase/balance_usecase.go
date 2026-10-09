package usecase

import (
	"context"
	"errors"
	"fmt"

	"nexus-video-games/internal/entity"
	"nexus-video-games/internal/gateway/email"
	"nexus-video-games/internal/gateway/payment"
	"nexus-video-games/internal/model"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BalanceUsecase struct {
	db          *gorm.DB
	balanceRepo *postgresql.BalanceRepository
	topupRepo   *postgresql.TopupRepository
	userRepo    *postgresql.UserRepository
	xenditGW    *payment.XenditGateway
	emailGW     *email.ResendGateway
}

func NewBalanceUsecase(
	db *gorm.DB,
	balanceRepo *postgresql.BalanceRepository,
	topupRepo *postgresql.TopupRepository,
	userRepo *postgresql.UserRepository,
	xenditGW *payment.XenditGateway,
	emailGW *email.ResendGateway,
) *BalanceUsecase {
	return &BalanceUsecase{
		db:          db,
		balanceRepo: balanceRepo,
		topupRepo:   topupRepo,
		userRepo:    userRepo,
		xenditGW:    xenditGW,
		emailGW:     emailGW,
	}
}

func (u *BalanceUsecase) GetBalance(ctx context.Context, userID uuid.UUID) (*model.BalanceResponse, error) {
	balance, err := u.balanceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.BalanceResponse{
		UserID:        balance.UserID,
		BalanceAmount: balance.BalanceAmount,
		Currency:      balance.Currency,
		UpdatedAt:     balance.UpdatedAt,
	}, nil
}

func (u *BalanceUsecase) CreateTopup(ctx context.Context, userID uuid.UUID, email string, req *model.TopupBalanceRequest) (*model.TopupResponse, error) {
	externalID := fmt.Sprintf("TOPUP-%s-%d", userID.String()[:8], uuid.New().ID())

	inv, err := u.xenditGW.CreateInvoice(ctx, payment.XenditInvoiceRequest{
		ExternalID:      externalID,
		Amount:          req.Amount,
		PayerEmail:      email,
		Description:     "Nexus Video Games Wallet Top-Up",
		InvoiceDuration: 86400,
	})
	if err != nil {
		return nil, err
	}

	topup := &entity.XenditTopup{
		UserID:              userID,
		XenditInvoiceID:     inv.ID,
		Amount:              req.Amount,
		Status:              entity.PaymentStatusPending,
		ExternalCheckoutURL: inv.InvoiceURL,
	}

	if err := u.topupRepo.Create(ctx, topup); err != nil {
		return nil, err
	}

	return &model.TopupResponse{
		InvoiceID:   topup.XenditInvoiceID,
		CheckoutURL: topup.ExternalCheckoutURL,
		Amount:      topup.Amount,
		Status:      string(topup.Status),
		CreatedAt:   topup.CreatedAt,
	}, nil
}

func (u *BalanceUsecase) HandleXenditWebhook(ctx context.Context, webhookToken string, payload *model.XenditInvoiceWebhookPayload) error {
	if !u.xenditGW.VerifyWebhookToken(webhookToken) {
		return errors.New("invalid webhook signature token")
	}

	topup, err := u.topupRepo.FindByInvoiceID(ctx, payload.ID)
	if err != nil {
		return err
	}

	if topup.Status == entity.PaymentStatusPaid {
		return nil // Idempotent execution
	}

	switch payload.Status {
	case "PAID", "SETTLED":
		err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := u.topupRepo.UpdateStatus(ctx, tx, topup.XenditInvoiceID, entity.PaymentStatusPaid); err != nil {
				return err
			}
			return u.balanceRepo.AddBalance(ctx, tx, topup.UserID, topup.Amount)
		})
		if err != nil {
			return err
		}

		user, _ := u.userRepo.FindByID(ctx, topup.UserID)
		if user != nil {
			_ = u.emailGW.SendTopupReceiptEmail(ctx, user.Email, topup.Amount, topup.XenditInvoiceID)
		}
	case "EXPIRED":
		_ = u.topupRepo.UpdateStatus(ctx, nil, topup.XenditInvoiceID, entity.PaymentStatusExpired)
	}

	return nil
}
