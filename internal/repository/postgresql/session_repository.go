package postgresql

import (
	"context"
	"time"

	"nexus-video-games/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// FindUnexpiredSession checks if the user is currently locked into an active session
func (r *SessionRepository) FindUnexpiredSession(ctx context.Context, userID uuid.UUID) (*entity.ActiveSession, error) {
	var session entity.ActiveSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpsertSession saves or refreshes the user's active playing session
func (r *SessionRepository) UpsertSession(ctx context.Context, session *entity.ActiveSession) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"games_owned_id", "external_game_id", "started_at", "expires_at"}),
	}).Create(session).Error
}

func (r *SessionRepository) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.ActiveSession{}).Error
}
