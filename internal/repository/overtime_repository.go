package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) domain.OvertimeRepository {
	return &overtimeRepository{db}
}

func (r *overtimeRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Overtime, error) {
	var overtimes []domain.Overtime
	query := r.db.WithContext(ctx).Model(&domain.Overtime{})

	// Explicit filtering (Senior BE approach)
	if empID, ok := filter["admin_id"]; ok {
		query = query.Where("overtimes.admin_id = ?", empID)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("overtimes.status = ?", status)
	}
	if last7, ok := filter["last_7_days"]; ok && last7 == true {
		query = query.Where("overtimes.created_at >= CURRENT_DATE - interval '7 days'")
	}
	if monthYear, ok := filter["month_year"]; ok {
		query = query.Where("overtimes.overtime_date::text LIKE ?", monthYear)
	}

	err := query.Select("overtimes.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = overtimes.admin_id").
		Order("overtimes.created_at DESC").
		Find(&overtimes).Error

	return overtimes, err
}

func (r *overtimeRepository) GetByID(ctx context.Context, id int) (*domain.Overtime, error) {
	var overtime domain.Overtime
	err := r.db.WithContext(ctx).Model(&domain.Overtime{}).
		Select("overtimes.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = overtimes.admin_id").
		First(&overtime, id).Error
	if err != nil {
		return nil, err
	}
	return &overtime, nil
}

func (r *overtimeRepository) Store(ctx context.Context, overtime *domain.Overtime) error {
	return r.db.WithContext(ctx).Create(overtime).Error
}

func (r *overtimeRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).Model(&domain.Overtime{}).Where("id = ?", id).Update("status", status).Error
}
