package model

import (
	"time"

	"github.com/google/uuid"
)

type LaunchGameResponse struct {
	SessionID   uuid.UUID `json:"session_id"`
	GameTitle   string    `json:"game_title"`
	RedirectURL string    `json:"redirect_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	Message     string    `json:"message"`
}

type ActiveSessionResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	GamesOwnedID   uuid.UUID `json:"games_owned_id"`
	ExternalGameID string    `json:"external_game_id"`
	StartedAt      time.Time `json:"started_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	SecondsLeft    int64     `json:"seconds_left"`
}
