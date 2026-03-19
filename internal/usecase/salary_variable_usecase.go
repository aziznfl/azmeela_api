package usecase

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type salaryVariableUsecase struct {
	svRepo domain.SalaryVariableRepository
}

func NewSalaryVariableUsecase(r domain.SalaryVariableRepository) domain.SalaryVariableUsecase {
	return &salaryVariableUsecase{
		svRepo: r,
	}
}

func (u *salaryVariableUsecase) Fetch(ctx context.Context) ([]domain.SalaryVariable, error) {
	return u.svRepo.Fetch(ctx)
}

func (u *salaryVariableUsecase) GetByID(ctx context.Context, id int) (*domain.SalaryVariable, error) {
	return u.svRepo.GetByID(ctx, id)
}

func (u *salaryVariableUsecase) Create(ctx context.Context, req *domain.SalaryVariableRequest) (*domain.SalaryVariable, error) {
	sv := &domain.SalaryVariable{
		Name:  req.Name,
		Type:  req.Type,
		Value: req.Value,
	}

	err := u.svRepo.Store(ctx, sv)
	if err != nil {
		return nil, err
	}
	return sv, nil
}

func (u *salaryVariableUsecase) Update(ctx context.Context, id int, req *domain.SalaryVariableRequest) (*domain.SalaryVariable, error) {
	sv, err := u.svRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sv.Name = req.Name
	sv.Type = req.Type
	sv.Value = req.Value

	err = u.svRepo.Update(ctx, sv)
	if err != nil {
		return nil, err
	}
	return sv, nil
}

func (u *salaryVariableUsecase) Delete(ctx context.Context, id int) error {
	return u.svRepo.Delete(ctx, id)
}
