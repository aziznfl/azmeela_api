package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) domain.LeaveRepository {
	return &leaveRepository{db}
}

func (r *leaveRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Leave, error) {
	var leaves []domain.Leave
	query := r.db.WithContext(ctx).Model(&domain.Leave{})

	// Sanitize and explicitly builder filters (Senior BE approach)
	if empID, ok := filter["admin_id"]; ok {
		query = query.Where("t_cuti.id_admin = ?", empID)
	}
	if types, ok := filter["type"]; ok {
		query = query.Where("t_cuti.grouping = ?", types)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("t_cuti.status = ?", status)
	}
	if upcoming, ok := filter["upcoming"]; ok && upcoming == true {
		query = query.Where("t_cuti.tanggal BETWEEN CURRENT_DATE AND (CURRENT_DATE + interval '7 days') AND t_cuti.status = 1")
	}
	if monthYear, ok := filter["month_year"]; ok {
		query = query.Where("t_cuti.tanggal::text LIKE ?", monthYear)
	}

	err := query.Select("t_cuti.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_cuti.id_admin").
		Order("t_cuti.tanggal DESC").
		Find(&leaves).Error

	return leaves, err
}

func (r *leaveRepository) GetByID(ctx context.Context, id int) (*domain.Leave, error) {
	var leave domain.Leave
	err := r.db.WithContext(ctx).Model(&domain.Leave{}).
		Select("t_cuti.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_cuti.id_admin").
		First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *leaveRepository) Store(ctx context.Context, leave *domain.Leave) error {
	return r.db.WithContext(ctx).Create(leave).Error
}

func (r *leaveRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).Model(&domain.Leave{}).Where("id_cuti = ?", id).Update("status", status).Error
}
