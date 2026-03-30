package domain

import "context"

// SalaryVariable represents configurable salary components (allowances, deductions, etc.)
type SalaryVariable struct {
	ID    int    `gorm:"primaryKey;column:id_variable;autoIncrement"`
	Name  string `gorm:"column:nama_variable"`
	Type  int    `gorm:"column:tipe_variable"`   // "allowance" or "deduction"
	Value int    `gorm:"column:nilai_variable"` // Fixed amount
}

func (SalaryVariable) TableName() string {
	return "t_variable_gaji"
}

// SalaryVariableRequest is the payload for creating/updating a salary variable
type SalaryVariableRequest struct {
	Name  string
	Type  int
	Value int
}

// SalaryVariableRepository represents the salary variable's repository contract
type SalaryVariableRepository interface {
	Fetch(ctx context.Context) ([]SalaryVariable, error)
	GetByID(ctx context.Context, id int) (*SalaryVariable, error)
	Store(ctx context.Context, sv *SalaryVariable) error
	Update(ctx context.Context, sv *SalaryVariable) error
	Delete(ctx context.Context, id int) error
}

// SalaryVariableUsecase represents the salary variable's usecase contract
type SalaryVariableUsecase interface {
	Fetch(ctx context.Context) ([]SalaryVariable, error)
	GetByID(ctx context.Context, id int) (*SalaryVariable, error)
	Create(ctx context.Context, req *SalaryVariableRequest) (*SalaryVariable, error)
	Update(ctx context.Context, id int, req *SalaryVariableRequest) (*SalaryVariable, error)
	Delete(ctx context.Context, id int) error
}
