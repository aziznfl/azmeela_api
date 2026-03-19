package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type adminTypeRepository struct {
	db *gorm.DB
}

func NewAdminTypeRepository(db *gorm.DB) domain.AdminTypeRepository {
	return &adminTypeRepository{db}
}

func (r *adminTypeRepository) Fetch(ctx context.Context) ([]domain.AdminType, error) {
	var types []domain.AdminType
	err := r.db.WithContext(ctx).Order("id ASC").Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (r *adminTypeRepository) GetByID(ctx context.Context, id int) (*domain.AdminType, error) {
	var at domain.AdminType
	err := r.db.WithContext(ctx).First(&at, id).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}
