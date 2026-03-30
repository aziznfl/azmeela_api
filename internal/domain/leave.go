package domain

import (
	"context"
	"time"
)

type Leave struct {
	ID           int       `gorm:"primaryKey;column:id_cuti;autoIncrement"`
	EmployeeID   int       `gorm:"column:id_admin"`
	Type         int       `gorm:"column:grouping"` // 0: leave, 1: sick leave
	LeaveDate    time.Time `gorm:"column:tanggal;type:date"`
	Durations    int       `gorm:"-"`
	Status       int       `gorm:"column:status"` // 0: created, 1: accepted, 2: rejected
	Description  string    `gorm:"column:keterangan"`
	EmployeeName string    `gorm:"-"`
}

func (Leave) TableName() string {
	return "t_cuti"
}

type LeaveResponse struct {
	ID           int
	EmployeeID   int
	EmployeeName string
	Type         int // 0: leave, 1: sick leave
	LeaveDate    time.Time
	Durations    int
	Status       int // 0: created, 1: accepted, 2: rejected
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LeaveRequest struct {
	Type        int
	LeaveDate   string
	Durations   int
	Description string
}

type LeaveStatusUpdate struct {
	Status int
}

type LeaveRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Leave, error)
	GetByID(ctx context.Context, id int) (*Leave, error)
	Store(ctx context.Context, leave *Leave) error
	UpdateStatus(ctx context.Context, id int, status int) error
}

type LeaveUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]LeaveResponse, error)
	RequestLeave(ctx context.Context, employeeID int, req *LeaveRequest) (*Leave, error)
	UpdateStatus(ctx context.Context, id int, req *LeaveStatusUpdate) error
}
