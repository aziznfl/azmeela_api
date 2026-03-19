package domain

import (
	"context"
	"time"
)

type Overtime struct {
	ID           int       `gorm:"primaryKey;column:id;autoIncrement"`
	EmployeeID   int       `gorm:"column:admin_id"`
	Date         time.Time `gorm:"column:overtime_date;type:date"`
	TimeIn       string    `gorm:"column:start_time;type:time"`
	TimeOut      string    `gorm:"column:end_time;type:time"`
	Status       int       `gorm:"column:status"` // 0: created, 1: approved, 2: disapproved
	Description  string    `gorm:"column:description"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	EmployeeName string    `gorm:"<-:false;column:employee_name"`
}

func (Overtime) TableName() string {
	return "overtimes"
}

type OvertimeResponse struct {
	ID           int
	EmployeeID   int
	EmployeeName string
	Date         time.Time
	TimeIn       string
	TimeOut      string
	Status       int // 0: created, 1: approved, 2: disapproved
	Description  string
}

type OvertimeRequest struct {
	Date        string
	TimeIn      string
	TimeOut     string
	Description string
}

// OvertimeStatusUpdate represents the payload to approve/reject an overtime request
type OvertimeStatusUpdate struct {
	Status int
}

type OvertimeRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Overtime, error)
	GetByID(ctx context.Context, id int) (*Overtime, error)
	Store(ctx context.Context, overtime *Overtime) error
	UpdateStatus(ctx context.Context, id int, status int) error
}

type OvertimeUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]OvertimeResponse, error)
	RequestOvertime(ctx context.Context, employeeID int, req *OvertimeRequest) (*Overtime, error)
	UpdateStatus(ctx context.Context, id int, req *OvertimeStatusUpdate) error
}
