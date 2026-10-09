package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email        string      `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string      `gorm:"type:varchar(255);not null" json:"-"`
	IsVerified   bool        `gorm:"type:boolean;default:false;not null" json:"is_verified"`
	CreatedAt    time.Time   `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Balance      UserBalance `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"balance,omitempty"`
}

func (User) TableName() string {
	return "users"
}
