package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type OvertimeResponse struct {
	ID           int       `json:"id"`
	EmployeeID   int       `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	Date         time.Time `json:"date"`
	TimeIn       string    `json:"time_in"`
	TimeOut      string    `json:"time_out"`
	Status       int       `json:"status"` // 0: created, 1: approved, 2: disapproved
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type OvertimeRequest struct {
	Date        string `json:"overtime_date" binding:"required"`
	TimeIn      string `json:"start_time" binding:"required"`
	TimeOut     string `json:"end_time" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type OvertimeStatusUpdate struct {
	Status int `json:"status" binding:"oneof=0 1 2"`
}

func ToOvertimeResponse(o *domain.Overtime) *OvertimeResponse {
	if o == nil {
		return nil
	}
	return &OvertimeResponse{
		ID:           o.ID,
		EmployeeID:   o.EmployeeID,
		EmployeeName: o.EmployeeName,
		Date:         o.Date,
		TimeIn:       o.TimeIn,
		TimeOut:      o.TimeOut,
		Status:       o.Status,
		Description:  o.Description,
	}
}

func ToOvertimeListResponse(items []domain.Overtime) []*OvertimeResponse {
	resps := make([]*OvertimeResponse, len(items))
	for i, item := range items {
		resps[i] = ToOvertimeResponse(&item)
	}
	return resps
}

func ToOvertimeResponseFromDomain(o *domain.OvertimeResponse) *OvertimeResponse {
	if o == nil {
		return nil
	}
	return &OvertimeResponse{
		ID:           o.ID,
		EmployeeID:   o.EmployeeID,
		EmployeeName: o.EmployeeName,
		Date:         o.Date,
		TimeIn:       o.TimeIn,
		TimeOut:      o.TimeOut,
		Status:       o.Status,
		Description:  o.Description,
	}
}

func ToOvertimeListResponseFromDomain(items []domain.OvertimeResponse) []*OvertimeResponse {
	resps := make([]*OvertimeResponse, len(items))
	for i, item := range items {
		resps[i] = ToOvertimeResponseFromDomain(&item)
	}
	return resps
}
