package domain

import (
	"context"
	"time"
)

// AdminType represents the admin_types table in the database
type AdminType struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(20);unique"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AdminType) TableName() string {
	return "admin_types"
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
