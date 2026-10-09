package entity

import (
	"time"

	"github.com/google/uuid"
)

type ActiveSession struct {
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	GamesOwnedID   uuid.UUID  `gorm:"type:uuid;not null" json:"games_owned_id"`
	ExternalGameID string     `gorm:"type:varchar(100);not null" json:"external_game_id"`
	StartedAt      time.Time  `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"started_at"`
	ExpiresAt      time.Time  `gorm:"type:timestamp with time zone;not null;index:idx_active_sessions_timeout,priority:2" json:"expires_at"`
	GameOwned      GamesOwned `gorm:"foreignKey:GamesOwnedID;references:ID" json:"game_owned,omitempty"`
}

func (ActiveSession) TableName() string {
	return "active_sessions"
}
