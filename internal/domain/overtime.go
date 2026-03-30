package domain

import (
	"context"
	"time"
)

type Overtime struct {
	ID           int       `gorm:"primaryKey;column:id_lembur;autoIncrement"`
	EmployeeID   int       `gorm:"column:id_admin"`
	Date         time.Time `gorm:"column:tanggal;type:date"`
	TimeIn       string    `gorm:"column:jam_masuk;type:time"`
	TimeOut      string    `gorm:"column:jam_keluar;type:time"`
	Status       int       `gorm:"column:status"` // 0: created, 1: approved, 2: disapproved
	Description  string    `gorm:"column:keterangan"`
	EmployeeName string    `gorm:"-"`
}

func (Overtime) TableName() string {
	return "t_lembur"
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
