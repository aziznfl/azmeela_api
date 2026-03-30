package domain

import (
	"context"
	"time"
)

type CashAdvance struct {
	ID           int       `gorm:"primaryKey;column:id_kasbon;autoIncrement"`
	EmployeeID   int       `gorm:"column:id_admin"`
	Amount       int       `gorm:"column:jumlah"`
	Purpose      string    `gorm:"column:keperluan"`
	Status       int       `gorm:"column:status"` // 0: created, 1: approved, 2: disapproved
	CreatedAt    time.Time `gorm:"column:tanggal;autoCreateTime"`
	EmployeeName string    `gorm:"-"`
}

func (CashAdvance) TableName() string {
	return "t_kasbon"
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
	ID         int       `gorm:"primaryKey;column:id_history;autoIncrement"`
	EmployeeID int       `gorm:"column:id_admin"`
	Date       time.Time `gorm:"column:tanggal;type:date"`
	Amount     int       `gorm:"column:jumlah"`
	Type       int       `gorm:"column:tipe"` // 1: debt, 2: payment
	Purpose    string    `gorm:"-"`
}

func (CashAdvanceHistory) TableName() string {
	return "t_kasbon_history"
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
