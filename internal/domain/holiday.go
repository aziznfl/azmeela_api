package domain

import (
	"context"
	"time"
)

type Holiday struct {
	ID          int       `gorm:"primaryKey;column:id;autoIncrement"`
	HolidayDate time.Time `gorm:"column:holiday_date;type:date"`
	Month       int       `gorm:"column:month"`
	Day         int       `gorm:"column:day"`
	Description string    `gorm:"column:description"`
	IsRecurring bool      `gorm:"column:is_recurring;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Holiday) TableName() string {
	return "holidays"
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
