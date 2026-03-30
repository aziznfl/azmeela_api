package domain

import (
	"context"
	"time"
)

// Employee represents the t_admin table in the database
type Employee struct {
	ID            int        `gorm:"primaryKey;column:id_admin;autoIncrement"`
	TypeID        int        `gorm:"column:id_admin_type"`
	Username      string     `gorm:"column:username;unique"`
	Password      string     `gorm:"column:password"` // Hashed password
	Name          string     `gorm:"column:nama_admin"`
	Active        int        `gorm:"column:status_admin"`
	Bio           string     `gorm:"column:bio"`
	BaseSalary    int        `gorm:"column:gaji_pokok"`
	ContractStart *time.Time `gorm:"column:kontrak_awal;type:date"`
	ContractEnd   *time.Time `gorm:"column:kontrak_akhir;type:date"`
	CV            string     `gorm:"column:cv"`
	Phone         string     `gorm:"column:no_hp"`
	DateSign      time.Time  `gorm:"column:date_sign;autoCreateTime"`
	Photo         string     `gorm:"column:photo"`
	AdminType     *AdminType `gorm:"foreignKey:TypeID"`
}

func (Employee) TableName() string {
	return "t_admin"
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
