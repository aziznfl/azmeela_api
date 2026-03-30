package domain

import (
	"context"
	"time"
)

// Attendance represents the presences table in the database
type Attendance struct {
	ID           int       `gorm:"primaryKey;column:id_presensi;autoIncrement"`
	EmployeeID   int       `gorm:"column:id_admin"`
	Date         time.Time `gorm:"column:tanggal;type:date"`
	TimeIn       string    `gorm:"column:jam_masuk;type:time"`
	TimeOut      *string   `gorm:"column:jam_keluar;type:time"`
	Location     *string   `gorm:"-"`
	Note         *string   `gorm:"-"`
	EmployeeName string    `gorm:"-"`
}

// TableName overrides the table name used by Gorm
func (Attendance) TableName() string {
	return "t_presensi"
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
