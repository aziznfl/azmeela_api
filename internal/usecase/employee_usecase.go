package usecase

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type employeeUsecase struct {
	employeeRepo domain.EmployeeRepository
}

// NewEmployeeUsecase will create new an employeeUsecase object representation of domain.EmployeeUsecase interface
func NewEmployeeUsecase(e domain.EmployeeRepository) domain.EmployeeUsecase {
	return &employeeUsecase{
		employeeRepo: e,
	}
}

func (u *employeeUsecase) Fetch(ctx context.Context) ([]domain.Employee, error) {
	return u.employeeRepo.Fetch(ctx)
}

func (u *employeeUsecase) GetByID(ctx context.Context, id int) (*domain.Employee, error) {
	return u.employeeRepo.GetByID(ctx, id)
}

func (u *employeeUsecase) Store(ctx context.Context, emp *domain.Employee) error {
	// In a complete app, we could add business logic here (e.g., hash password, validate bio...)
	return u.employeeRepo.Store(ctx, emp)
}

func (u *employeeUsecase) Update(ctx context.Context, emp *domain.Employee) error {
	return u.employeeRepo.Update(ctx, emp)
}

func (u *employeeUsecase) Delete(ctx context.Context, id int) error {
	return u.employeeRepo.Delete(ctx, id)
}
