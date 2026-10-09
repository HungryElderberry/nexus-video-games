package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserBalance struct {
	UserID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	BalanceAmount int64     `gorm:"type:bigint;default:0;not null" json:"balance_amount"`
	Currency      string    `gorm:"type:varchar(3);default:'IDR';not null" json:"currency"`
	UpdatedAt     time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (UserBalance) TableName() string {
	return "user_balances"
}
