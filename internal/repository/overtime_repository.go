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
		query = query.Where("t_lembur.id_admin = ?", empID)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("t_lembur.status = ?", status)
	}
	if monthYear, ok := filter["month_year"]; ok {
		query = query.Where("t_lembur.tanggal::text LIKE ?", monthYear)
	}

	err := query.Select("t_lembur.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_lembur.id_admin").
		Order("t_lembur.tanggal DESC").
		Find(&overtimes).Error

	return overtimes, err
}

func (r *overtimeRepository) GetByID(ctx context.Context, id int) (*domain.Overtime, error) {
	var overtime domain.Overtime
	err := r.db.WithContext(ctx).Model(&domain.Overtime{}).
		Select("t_lembur.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_lembur.id_admin").
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
	return r.db.WithContext(ctx).Model(&domain.Overtime{}).Where("id_lembur = ?", id).Update("status", status).Error
}
