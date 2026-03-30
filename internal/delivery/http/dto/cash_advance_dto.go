package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type CashAdvanceResponse struct {
	ID           int       `json:"id"`
	EmployeeID   int       `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	Amount       int       `json:"amount"`
	Purpose      string    `json:"purpose"`
	Status       int       `json:"status"` // 0: created, 1: approved, 2: disapproved
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CashAdvanceRequest struct {
	Amount  int    `json:"amount" binding:"required,gt=0"`
	Purpose string `json:"purpose" binding:"required"`
}

type CashAdvanceStatusUpdate struct {
	Status int `json:"status" binding:"oneof=0 1 2"`
}

type CashAdvancePayment struct {
	EmployeeID int       `json:"employee_id" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	Amount     int       `json:"amount" binding:"required,gt=0"`
}

func ToCashAdvanceResponse(ca *domain.CashAdvance) *CashAdvanceResponse {
	if ca == nil {
		return nil
	}
	return &CashAdvanceResponse{
		ID:           ca.ID,
		EmployeeID:   ca.EmployeeID,
		EmployeeName: ca.EmployeeName,
		Amount:       ca.Amount,
		Purpose:      ca.Purpose,
		Status:       ca.Status,
		CreatedAt:    ca.CreatedAt,
	}
}

func ToCashAdvanceListResponse(items []domain.CashAdvance) []*CashAdvanceResponse {
	resps := make([]*CashAdvanceResponse, len(items))
	for i, item := range items {
		resps[i] = ToCashAdvanceResponse(&item)
	}
	return resps
}

func ToCashAdvanceResponseFromDomain(a *domain.CashAdvanceResponse) *CashAdvanceResponse {
	if a == nil {
		return nil
	}
	return &CashAdvanceResponse{
		ID:           a.ID,
		EmployeeID:   a.EmployeeID,
		EmployeeName: a.EmployeeName,
		Amount:       a.Amount,
		Purpose:      a.Purpose,
		Status:       a.Status,
		CreatedAt:    a.CreatedAt,
	}
}

func ToCashAdvanceListResponseFromDomain(advances []domain.CashAdvanceResponse) []*CashAdvanceResponse {
	resps := make([]*CashAdvanceResponse, len(advances))
	for i, ca := range advances {
		resps[i] = ToCashAdvanceResponseFromDomain(&ca)
	}
	return resps
}
