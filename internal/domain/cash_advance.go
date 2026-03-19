package domain

import (
	"context"
	"time"
)

type CashAdvance struct {
	ID           int       `gorm:"primaryKey;column:id;autoIncrement"`
	EmployeeID   int       `gorm:"column:admin_id"`
	Amount       int       `gorm:"column:amount"`
	Purpose      string    `gorm:"column:purpose"`
	Status       int       `gorm:"column:status"` // 0: created, 1: approved, 2: disapproved
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	EmployeeName string    `gorm:"<-:false;column:employee_name"`
}

func (CashAdvance) TableName() string {
	return "cash_advances"
}

type CashAdvanceResponse struct {
	ID           int
	EmployeeID   int
	EmployeeName string
	Amount       int
	Purpose      string
	Status       int // 0: created, 1: approved, 2: disapproved
	CreatedAt    time.Time
}

type CashAdvanceHistory struct {
	ID         int       `gorm:"primaryKey;column:id;autoIncrement"`
	EmployeeID int       `gorm:"column:admin_id"`
	Date       time.Time `gorm:"column:tanggal;type:date"`
	Amount     int       `gorm:"column:amount"`
	Type       int       `gorm:"column:tipe"` // 1: debt, 2: payment
	Purpose    string    `gorm:"column:purpose"`
}

func (CashAdvanceHistory) TableName() string {
	return "cash_advance_histories"
}

type CashAdvanceRequest struct {
	Amount  int
	Purpose string
}

type CashAdvanceStatusUpdate struct {
	Status int
}

type CashAdvancePayment struct {
	EmployeeID int
	Date       time.Time
	Amount     int
}

type CashAdvanceRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]CashAdvance, error)
	GetByID(ctx context.Context, id int) (*CashAdvance, error)
	Store(ctx context.Context, ca *CashAdvance) error
	UpdateStatus(ctx context.Context, id int, status int) error
	StoreHistory(ctx context.Context, history *CashAdvanceHistory) error
}

type CashAdvanceUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}) ([]CashAdvanceResponse, error)
	RequestCashAdvance(ctx context.Context, employeeID int, req *CashAdvanceRequest) (*CashAdvance, error)
	UpdateStatus(ctx context.Context, id int, req *CashAdvanceStatusUpdate) error
	AddPayment(ctx context.Context, req *CashAdvancePayment) error
}
