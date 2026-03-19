package usecase

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type adminTypeUsecase struct {
	repo domain.AdminTypeRepository
}

func NewAdminTypeUsecase(repo domain.AdminTypeRepository) domain.AdminTypeUsecase {
	return &adminTypeUsecase{
		repo: repo,
	}
}

func (u *adminTypeUsecase) Fetch(ctx context.Context) ([]domain.AdminType, error) {
	return u.repo.Fetch(ctx)
}
