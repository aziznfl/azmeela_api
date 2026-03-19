package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type cashAdvanceRepository struct {
	db *gorm.DB
}

func NewCashAdvanceRepository(db *gorm.DB) domain.CashAdvanceRepository {
	return &cashAdvanceRepository{db}
}

func (r *cashAdvanceRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.CashAdvance, error) {
	var results []domain.CashAdvance
	query := r.db.WithContext(ctx).Model(&domain.CashAdvance{})

	// Sanitize and builder filters (Senior BE approach)
	if empID, ok := filter["admin_id"]; ok {
		query = query.Where("cash_advances.admin_id = ?", empID)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("cash_advances.status = ?", status)
	}
	if last7, ok := filter["last_7_days"]; ok && last7 == true {
		query = query.Where("cash_advances.created_at >= CURRENT_DATE - interval '7 days'")
	}
	if monthYear, ok := filter["month_year"]; ok {
		query = query.Where("cash_advances.created_at::text LIKE ?", monthYear)
	}

	err := query.Select("cash_advances.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = cash_advances.admin_id").
		Order("cash_advances.created_at DESC").
		Find(&results).Error

	return results, err
}

func (r *cashAdvanceRepository) GetByID(ctx context.Context, id int) (*domain.CashAdvance, error) {
	var ca domain.CashAdvance
	err := r.db.WithContext(ctx).Model(&domain.CashAdvance{}).
		Select("cash_advances.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = cash_advances.admin_id").
		First(&ca, id).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func (r *cashAdvanceRepository) Store(ctx context.Context, ca *domain.CashAdvance) error {
	return r.db.WithContext(ctx).Create(ca).Error
}

func (r *cashAdvanceRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).Model(&domain.CashAdvance{}).Where("id = ?", id).Update("status", status).Error
}

func (r *cashAdvanceRepository) StoreHistory(ctx context.Context, history *domain.CashAdvanceHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}
