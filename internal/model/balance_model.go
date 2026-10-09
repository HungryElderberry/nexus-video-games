package model

import (
	"time"

	"github.com/google/uuid"
)

type TopupBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required,min=10000"` // Minimum Topup IDR 10k
}

type TopupResponse struct {
	InvoiceID   string    `json:"invoice_id"`
	CheckoutURL string    `json:"checkout_url"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type BalanceResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	BalanceAmount int64     `json:"balance_amount"`
	Currency      string    `json:"currency"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type XenditInvoiceWebhookPayload struct {
	ID            string `json:"id"`
	ExternalID    string `json:"external_id"`
	UserID        string `json:"user_id"`
	Status        string `json:"status"`
	PaidAmount    int64  `json:"paid_amount"`
	PayerEmail    string `json:"payer_email"`
	PaymentMethod string `json:"payment_method"`
}
