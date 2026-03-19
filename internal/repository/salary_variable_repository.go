package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type salaryVariableRepository struct {
	db *gorm.DB
}

func NewSalaryVariableRepository(db *gorm.DB) domain.SalaryVariableRepository {
	return &salaryVariableRepository{db}
}

func (r *salaryVariableRepository) Fetch(ctx context.Context) ([]domain.SalaryVariable, error) {
	var variables []domain.SalaryVariable
	err := r.db.WithContext(ctx).Order("type ASC, name ASC").Find(&variables).Error
	if err != nil {
		return nil, err
	}
	return variables, nil
}

func (r *salaryVariableRepository) GetByID(ctx context.Context, id int) (*domain.SalaryVariable, error) {
	var sv domain.SalaryVariable
	err := r.db.WithContext(ctx).First(&sv, id).Error
	if err != nil {
		return nil, err
	}
	return &sv, nil
}

func (r *salaryVariableRepository) Store(ctx context.Context, sv *domain.SalaryVariable) error {
	return r.db.WithContext(ctx).Create(sv).Error
}

func (r *salaryVariableRepository) Update(ctx context.Context, sv *domain.SalaryVariable) error {
	return r.db.WithContext(ctx).Save(sv).Error
}

func (r *salaryVariableRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.SalaryVariable{}, id).Error
}
