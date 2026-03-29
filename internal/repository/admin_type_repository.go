package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type adminTypeRepository struct {
	BaseRepository[domain.AdminType]
}

func NewAdminTypeRepository(db *gorm.DB) domain.AdminTypeRepository {
	return &adminTypeRepository{
		BaseRepository: BaseRepository[domain.AdminType]{db: db},
	}
}

func (r *adminTypeRepository) Fetch(ctx context.Context) ([]domain.AdminType, error) {
	res, _, err := r.BaseRepository.Fetch(ctx, nil, 0, 100)
	return res, err
}
