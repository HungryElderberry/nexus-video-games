package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusExpired PaymentStatus = "EXPIRED"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type XenditTopup struct {
	ID                  uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID              uuid.UUID     `gorm:"type:uuid;not null;index:idx_xendit_topups_user" json:"user_id"`
	XenditInvoiceID     string        `gorm:"type:varchar(255);unique;not null" json:"xendit_invoice_id"`
	Amount              int64         `gorm:"type:bigint;not null" json:"amount"`
	Status              PaymentStatus `gorm:"type:payment_status;default:'PENDING';not null" json:"status"`
	ExternalCheckoutURL string        `gorm:"type:text" json:"external_checkout_url"`
	CreatedAt           time.Time     `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time     `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (XenditTopup) TableName() string {
	return "xendit_topups"
}
