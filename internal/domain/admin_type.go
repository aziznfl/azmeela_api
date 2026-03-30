package domain

import (
	"context"
)

// AdminType represents the admin_types table in the database
type AdminType struct {
	ID   int    `gorm:"primaryKey;column:id_admin_type;autoIncrement"`
	Name string `gorm:"column:nama_admin_type;type:varchar(20);unique"`
}

func (AdminType) TableName() string {
	return "t_admin_type"
}

// AdminTypeRepository represent the admin type's repository contract
type AdminTypeRepository interface {
	Fetch(ctx context.Context) ([]AdminType, error)
	GetByID(ctx context.Context, id int) (*AdminType, error)
}

// AdminTypeUsecase represent the admin type's usecase contract
type AdminTypeUsecase interface {
	Fetch(ctx context.Context) ([]AdminType, error)
}
