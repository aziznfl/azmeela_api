package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type holidayRepository struct {
	db *gorm.DB
}

func NewHolidayRepository(db *gorm.DB) domain.HolidayRepository {
	return &holidayRepository{db}
}

func (r *holidayRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Holiday, error) {
	var holidays []domain.Holiday
	query := r.db.WithContext(ctx)

	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}

	err := query.Find(&holidays).Error
	if err != nil {
		return nil, err
	}
	return holidays, nil
}

func (r *holidayRepository) GetByID(ctx context.Context, id int) (*domain.Holiday, error) {
	var holiday domain.Holiday
	err := r.db.WithContext(ctx).First(&holiday, id).Error
	if err != nil {
		return nil, err
	}
	return &holiday, nil
}

func (r *holidayRepository) Store(ctx context.Context, holiday *domain.Holiday) error {
	return r.db.WithContext(ctx).Create(holiday).Error
}

func (r *holidayRepository) Update(ctx context.Context, holiday *domain.Holiday) error {
	return r.db.WithContext(ctx).Save(holiday).Error
}

func (r *holidayRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Holiday{}, id).Error
}
