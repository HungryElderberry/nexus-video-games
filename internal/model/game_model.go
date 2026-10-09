package model

import (
	"time"

	"github.com/google/uuid"
)

type GameCatalogItemResponse struct {
	ExternalGameID string `json:"external_game_id"`
	Title          string `json:"title"`
	PriceInBalance int64  `json:"price_in_balance"`
	IsAvailable    bool   `json:"is_available"`
}

type BuyGameRequest struct {
	ExternalGameID string `json:"external_game_id" binding:"required"`
}

type UserGameLibraryResponse struct {
	ID              uuid.UUID `json:"id"`
	ExternalGameID  string    `json:"external_game_id"`
	Title           string    `json:"title"`
	AcquiredAt      time.Time `json:"acquired_at"`
	PlaytimeMinutes int       `json:"playtime_minutes"`
}
