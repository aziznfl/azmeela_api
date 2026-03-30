package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type LeaveResponse struct {
	ID           int       `json:"id"`
	EmployeeID   int       `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	Type         int       `json:"type"` // 0: leave, 1: sick leave
	LeaveDate    time.Time `json:"leave_date"`
	Durations    int       `json:"durations"`
	Status       int       `json:"status"` // 0: created, 1: accepted, 2: rejected
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LeaveRequest struct {
	Type        int    `json:"type" binding:"oneof=0 1"` // 0 for leave, 1 for sick leave
	LeaveDate   string `json:"leave_date" binding:"required"`
	Durations   int    `json:"durations" binding:"required,min=1"`
	Description string `json:"description" binding:"required"`
}

type LeaveStatusUpdate struct {
	Status int `json:"status" binding:"oneof=0 1 2"`
}

func ToLeaveResponse(l *domain.Leave) *LeaveResponse {
	if l == nil {
		return nil
	}
	return &LeaveResponse{
		ID:           l.ID,
		EmployeeID:   l.EmployeeID,
		EmployeeName: l.EmployeeName,
		Type:         l.Type,
		LeaveDate:    l.LeaveDate,
		Durations:    l.Durations,
		Status:       l.Status,
		Description:  l.Description,
	}
}

func ToLeaveListResponse(items []domain.Leave) []*LeaveResponse {
	resps := make([]*LeaveResponse, len(items))
	for i, item := range items {
		resps[i] = ToLeaveResponse(&item)
	}
	return resps
}

func ToLeaveResponseFromDomain(l *domain.LeaveResponse) *LeaveResponse {
	if l == nil {
		return nil
	}
	return &LeaveResponse{
		ID:           l.ID,
		EmployeeID:   l.EmployeeID,
		EmployeeName: l.EmployeeName,
		Type:         l.Type,
		LeaveDate:    l.LeaveDate,
		Durations:    l.Durations,
		Status:       l.Status,
		Description:  l.Description,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func ToLeaveListResponseFromDomain(leaves []domain.LeaveResponse) []*LeaveResponse {
	resps := make([]*LeaveResponse, len(leaves))
	for i, l := range leaves {
		resps[i] = ToLeaveResponseFromDomain(&l)
	}
	return resps
}
