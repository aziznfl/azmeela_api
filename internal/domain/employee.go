package domain

import (
	"context"
	"time"
)

// Employee represents the t_admin table in the database
type Employee struct {
	ID            int        `gorm:"primaryKey;column:id;autoIncrement"`
	TypeID        int        `gorm:"column:admin_type_id"`
	Username      string     `gorm:"column:username;unique"`
	Password      string     `gorm:"column:password"` // Hashed password
	Name          string     `gorm:"column:name"`
	Active        bool       `gorm:"column:active"`
	Bio           string     `gorm:"column:bio"`
	BaseSalary    int        `gorm:"column:base_salary"`
	ContractStart *time.Time `gorm:"column:contract_start;type:date"`
	ContractEnd   *time.Time `gorm:"column:contract_end;type:date"`
	CV            string     `gorm:"column:cv"`
	Phone         string     `gorm:"column:phone"`
	AdminType     *AdminType `gorm:"foreignKey:TypeID"`
}

func (Employee) TableName() string {
	return "admins"
}

// EmployeeRepository represent the employee's repository contract
type EmployeeRepository interface {
	Fetch(ctx context.Context) ([]Employee, error)
	GetByID(ctx context.Context, id int) (*Employee, error)
	GetByUsername(ctx context.Context, username string) (*Employee, error)
	Store(ctx context.Context, emp *Employee) error
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id int) error
}

// EmployeeUsecase represent the employee's usecases
type EmployeeUsecase interface {
	Fetch(ctx context.Context) ([]Employee, error)
	GetByID(ctx context.Context, id int) (*Employee, error)
	Store(ctx context.Context, emp *Employee) error
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id int) error
}
