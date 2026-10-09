package postgresql

import (
	"context"

	"nexus-video-games/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// UpsertCatalog inserts or updates catalog records synced from IsThereAnyDeal API
func (r *GameRepository) UpsertCatalog(ctx context.Context, games []entity.GamesCatalog) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_game_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "price_in_balance", "is_available"}),
	}).Create(&games).Error
}

func (r *GameRepository) FindAvailableCatalog(ctx context.Context) ([]entity.GamesCatalog, error) {
	var games []entity.GamesCatalog
	err := r.db.WithContext(ctx).Where("is_available = ?", true).Find(&games).Error
	return games, err
}

func (r *GameRepository) FindCatalogByID(ctx context.Context, externalID string) (*entity.GamesCatalog, error) {
	var game entity.GamesCatalog
	err := r.db.WithContext(ctx).Where("external_game_id = ?", externalID).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *GameRepository) CreateOwned(ctx context.Context, tx *gorm.DB, owned *entity.GamesOwned) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(owned).Error
}

func (r *GameRepository) FindOwnedGame(ctx context.Context, userID uuid.UUID, externalGameID string) (*entity.GamesOwned, error) {
	var owned entity.GamesOwned
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND external_game_id = ?", userID, externalGameID).
		First(&owned).Error
	if err != nil {
		return nil, err
	}
	return &owned, nil
}

func (r *GameRepository) FindUserLibrary(ctx context.Context, userID uuid.UUID) ([]entity.GamesOwned, error) {
	var library []entity.GamesOwned
	err := r.db.WithContext(ctx).
		Preload("Game").
		Where("user_id = ?", userID).
		Find(&library).Error
	return library, err
}

func (r *GameRepository) CountOwnedGames(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.GamesOwned{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *GameRepository) IncrementPlaytime(ctx context.Context, id uuid.UUID, minutes int) error {
	return r.db.WithContext(ctx).Model(&entity.GamesOwned{}).
		Where("id = ?", id).
		Update("playtime_minutes", gorm.Expr("playtime_minutes + ?", minutes)).Error
}
