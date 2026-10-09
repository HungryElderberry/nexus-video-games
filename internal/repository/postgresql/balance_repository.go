package postgresql

import (
	"context"

	"nexus-video-games/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) Create(ctx context.Context, balance *entity.UserBalance) error {
	return r.db.WithContext(ctx).Create(balance).Error
}

func (r *BalanceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserBalance, error) {
	var balance entity.UserBalance
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *BalanceRepository) AddBalance(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int64) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	var balance entity.UserBalance
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).First(&balance).Error; err != nil {
		return err
	}

	return db.WithContext(ctx).Model(&entity.UserBalance{}).
		Where("user_id = ?", userID).
		Update("balance_amount", gorm.Expr("balance_amount + ?", amount)).Error
}

func (r *BalanceRepository) DeductBalance(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int64) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Model(&entity.UserBalance{}).
		Where("user_id = ? AND balance_amount >= ?", userID, amount).
		Update("balance_amount", gorm.Expr("balance_amount - ?", amount))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Insufficient funds or user missing
	}
	return nil
}
