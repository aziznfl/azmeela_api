package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository will create an object that represent the domain.EmployeeRepository interface
func NewEmployeeRepository(db *gorm.DB) domain.EmployeeRepository {
	return &employeeRepository{db}
}

func (r *employeeRepository) Fetch(ctx context.Context) ([]domain.Employee, error) {
	var employees []domain.Employee
	err := r.db.WithContext(ctx).
		Joins("AdminType").
		Order("admins.active DESC").
		Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id int) (*domain.Employee, error) {
	var employee domain.Employee
	err := r.db.WithContext(ctx).Joins("AdminType").First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByUsername(ctx context.Context, username string) (*domain.Employee, error) {
	var employee domain.Employee
	err := r.db.WithContext(ctx).Joins("AdminType").Where("username = ?", username).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) Store(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Updates(emp).Error
}

func (r *employeeRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Employee{}, id).Error
}
