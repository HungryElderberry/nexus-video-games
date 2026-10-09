package postgresql

import (
	"context"

	"nexus-video-games/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopupRepository struct {
	db *gorm.DB
}

func NewTopupRepository(db *gorm.DB) *TopupRepository {
	return &TopupRepository{db: db}
}

func (r *TopupRepository) Create(ctx context.Context, topup *entity.XenditTopup) error {
	return r.db.WithContext(ctx).Create(topup).Error
}

func (r *TopupRepository) FindByInvoiceID(ctx context.Context, invoiceID string) (*entity.XenditTopup, error) {
	var topup entity.XenditTopup
	err := r.db.WithContext(ctx).Where("xendit_invoice_id = ?", invoiceID).First(&topup).Error
	if err != nil {
		return nil, err
	}
	return &topup, nil
}

func (r *TopupRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, invoiceID string, status entity.PaymentStatus) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Model(&entity.XenditTopup{}).
		Where("xendit_invoice_id = ?", invoiceID).
		Update("status", status).Error
}

func (r *TopupRepository) GetTotalCashSpentByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.XenditTopup{}).
		Where("user_id = ? AND status = ?", userID, entity.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *TopupRepository) GetMonthlyCashSpentByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.XenditTopup{}).
		Where("user_id = ? AND status = ? AND created_at >= date_trunc('month', CURRENT_DATE)", userID, entity.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
