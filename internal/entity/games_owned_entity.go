package entity

import (
	"time"

	"github.com/google/uuid"
)

type GamesOwned struct {
	ID              uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID          uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:unique_user_game;index:idx_games_owned_user" json:"user_id"`
	ExternalGameID  string       `gorm:"type:varchar(100);not null;uniqueIndex:unique_user_game" json:"external_game_id"`
	AcquiredAt      time.Time    `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"acquired_at"`
	PlaytimeMinutes int          `gorm:"type:int;default:0;not null" json:"playtime_minutes"`
	Game            GamesCatalog `gorm:"foreignKey:ExternalGameID;references:ExternalGameID" json:"game,omitempty"`
}

func (GamesOwned) TableName() string {
	return "games_owned"
}
