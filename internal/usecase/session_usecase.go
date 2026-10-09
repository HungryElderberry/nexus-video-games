package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nexus-video-games/internal/entity"
	"nexus-video-games/internal/model"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type SessionUsecase struct {
	sessionRepo *postgresql.SessionRepository
	gameRepo    *postgresql.GameRepository
	appURL      string
}

func NewSessionUsecase(
	sessionRepo *postgresql.SessionRepository,
	gameRepo *postgresql.GameRepository,
	viper *viper.Viper,
) *SessionUsecase {
	appPort := viper.GetString("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &SessionUsecase{
		sessionRepo: sessionRepo,
		gameRepo:    gameRepo,
		appURL:      fmt.Sprintf("http://localhost:%s", appPort),
	}
}

func (u *SessionUsecase) LaunchGame(ctx context.Context, userID uuid.UUID, externalGameID string) (*model.LaunchGameResponse, error) {
	owned, err := u.gameRepo.FindOwnedGame(ctx, userID, externalGameID)
	if err != nil {
		return nil, errors.New("you do not own this game")
	}

	// Enforce 1 active game at a time rule via unexpired active sessions
	active, _ := u.sessionRepo.FindUnexpiredSession(ctx, userID)
	if active != nil {
		secondsLeft := int64(time.Until(active.ExpiresAt).Seconds())
		return nil, fmt.Errorf("active session already running for game '%s'. Remaining cooldown: %d seconds", active.ExternalGameID, secondsLeft)
	}

	expiresAt := time.Now().Add(60 * time.Second) // 60 seconds auto-expiry lock window
	session := &entity.ActiveSession{
		UserID:         userID,
		GamesOwnedID:   owned.ID,
		ExternalGameID: owned.ExternalGameID,
		StartedAt:      time.Now(),
		ExpiresAt:      expiresAt,
	}

	if err := u.sessionRepo.UpsertSession(ctx, session); err != nil {
		return nil, err
	}

	_ = u.gameRepo.IncrementPlaytime(ctx, owned.ID, 1)

	catalogGame, _ := u.gameRepo.FindCatalogByID(ctx, externalGameID)
	title := externalGameID
	if catalogGame != nil {
		title = catalogGame.Title
	}

	redirectURI := fmt.Sprintf("%s/games/%s", u.appURL, externalGameID)

	return &model.LaunchGameResponse{
		SessionID:   owned.ID,
		GameTitle:   title,
		RedirectURL: redirectURI,
		ExpiresAt:   expiresAt,
		Message:     "Session initialized. Redirecting to game route.",
	}, nil
}

func (u *SessionUsecase) GetActiveSession(ctx context.Context, userID uuid.UUID) (*model.ActiveSessionResponse, error) {
	session, err := u.sessionRepo.FindUnexpiredSession(ctx, userID)
	if err != nil {
		return nil, errors.New("no active game session currently running")
	}

	return &model.ActiveSessionResponse{
		UserID:         session.UserID,
		GamesOwnedID:   session.GamesOwnedID,
		ExternalGameID: session.ExternalGameID,
		StartedAt:      session.StartedAt,
		ExpiresAt:      session.ExpiresAt,
		SecondsLeft:    int64(time.Until(session.ExpiresAt).Seconds()),
	}, nil
}
