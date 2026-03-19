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
		query = query.Where("leaves.admin_id = ?", empID)
	}
	if types, ok := filter["type"]; ok {
		query = query.Where("leaves.type = ?", types)
	}
	if status, ok := filter["status"]; ok {
		query = query.Where("leaves.status = ?", status)
	}
	if upcoming, ok := filter["upcoming"]; ok && upcoming == true {
		query = query.Where("leaves.leave_date BETWEEN CURRENT_DATE AND (CURRENT_DATE + interval '7 days') AND leaves.status = 1")
	}
	if last7, ok := filter["last_7_days"]; ok && last7 == true {
		query = query.Where("leaves.created_at >= CURRENT_DATE - interval '7 days'")
	}
	if monthYear, ok := filter["month_year"]; ok {
		// Optimized date range query would be better, but keep current logic for now with safe execution
		query = query.Where("leaves.leave_date::text LIKE ?", monthYear)
	}

	err := query.Select("leaves.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = leaves.admin_id").
		Order("leaves.leave_date DESC").
		Find(&leaves).Error

	return leaves, err
}

func (r *leaveRepository) GetByID(ctx context.Context, id int) (*domain.Leave, error) {
	var leave domain.Leave
	err := r.db.WithContext(ctx).Model(&domain.Leave{}).
		Select("leaves.*, admins.name as employee_name").
		Joins("LEFT JOIN admins ON admins.id = leaves.admin_id").
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
	return r.db.WithContext(ctx).Model(&domain.Leave{}).Where("id = ?", id).Update("status", status).Error
}
