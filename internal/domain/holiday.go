package domain

import (
	"context"
	"time"
)

type Holiday struct {
	ID          int       `gorm:"primaryKey;column:id_libur;autoIncrement"`
	HolidayDate time.Time `gorm:"column:tanggal;type:date"`
	Month       int       `gorm:"-"`
	Day         int       `gorm:"-"`
	Description string    `gorm:"column:nama"`
	IsRecurring int       `gorm:"column:berulang;default:1"`
}

func (Holiday) TableName() string {
	return "t_libur"
}

type HolidayRequest struct {
	HolidayDate string
	Description string
	IsRecurring bool
}

type HolidayRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Holiday, error)
	GetByID(ctx context.Context, id int) (*Holiday, error)
	Store(ctx context.Context, holiday *Holiday) error
	Update(ctx context.Context, holiday *Holiday) error
	Delete(ctx context.Context, id int) error
}

type HolidayUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Holiday, error)
	Create(ctx context.Context, req *HolidayRequest) (*Holiday, error)
	Update(ctx context.Context, id int, req *HolidayRequest) (*Holiday, error)
	Delete(ctx context.Context, id int) error
}
