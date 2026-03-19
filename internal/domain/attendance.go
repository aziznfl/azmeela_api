package domain

import (
	"context"
	"time"
)

// Attendance represents the presences table in the database
type Attendance struct {
	ID           int       `gorm:"primaryKey;column:id;autoIncrement"`
	EmployeeID   int       `gorm:"column:employee_id"`
	Date         time.Time `gorm:"column:presence_date;type:date"`
	TimeIn       string    `gorm:"column:start_time;type:time"`
	TimeOut      *string   `gorm:"column:end_time;type:time"`
	Location     *string   `gorm:"column:location"`
	Note         *string   `gorm:"column:note"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	EmployeeName string    `gorm:"<-:false;column:employee_name"`
}

// TableName overrides the table name used by Gorm
func (Attendance) TableName() string {
	return "presences"
}

type AttendanceResponse struct {
	ID           int
	EmployeeID   int
	EmployeeName string
	Date         time.Time
	TimeIn       string
	TimeOut      *string
	Location     *string
	Note         *string
}

type AttendanceRequest struct {
	Location string
	Note     string
}

// AttendanceRepository exposes the database operations for Attendances
type AttendanceRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]Attendance, error)
	GetByDateAndEmployee(ctx context.Context, date time.Time, employeeID int) (*Attendance, error)
	GetTodayAttendances(ctx context.Context) ([]Attendance, error)
	Store(ctx context.Context, attendance *Attendance) error
	Update(ctx context.Context, attendance *Attendance) error
}

// AttendanceUsecase exposes the business logic operations for Attendances
type AttendanceUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]AttendanceResponse, error)
	ClockIn(ctx context.Context, employeeID int, req *AttendanceRequest) (*Attendance, error)
	ClockOut(ctx context.Context, employeeID int, req *AttendanceRequest) (*Attendance, error)
	GetTodayAttendances(ctx context.Context) ([]AttendanceResponse, error)
}
