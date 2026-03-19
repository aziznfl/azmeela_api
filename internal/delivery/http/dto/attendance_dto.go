package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type AttendanceResponse struct {
	ID           int       `json:"id"`
	EmployeeID   int       `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	Date         time.Time `json:"date"`
	TimeIn       string    `json:"time_in"`
	TimeOut      *string   `json:"time_out"`
	Location     *string   `json:"location"`
	Note         *string   `json:"note"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AttendanceRequest struct {
	Location string `json:"location"`
	Note     string `json:"note"`
}

func ToAttendanceResponse(a *domain.Attendance) *AttendanceResponse {
	if a == nil {
		return nil
	}
	return &AttendanceResponse{
		ID:           a.ID,
		EmployeeID:   a.EmployeeID,
		EmployeeName: a.EmployeeName,
		Date:         a.Date,
		TimeIn:       a.TimeIn,
		TimeOut:      a.TimeOut,
		Location:     a.Location,
		Note:         a.Note,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func ToAttendanceListResponse(attendances []domain.Attendance) []*AttendanceResponse {
	resps := make([]*AttendanceResponse, len(attendances))
	for i, a := range attendances {
		resps[i] = ToAttendanceResponse(&a)
	}
	return resps
}

func ToAttendanceResponseFromDomain(a *domain.AttendanceResponse) *AttendanceResponse {
	if a == nil {
		return nil
	}
	return &AttendanceResponse{
		ID:           a.ID,
		EmployeeID:   a.EmployeeID,
		EmployeeName: a.EmployeeName,
		Date:         a.Date,
		TimeIn:       a.TimeIn,
		TimeOut:      a.TimeOut,
		Location:     a.Location,
		Note:         a.Note,
	}
}

func ToAttendanceListResponseFromDomain(attendances []domain.AttendanceResponse) []*AttendanceResponse {
	resps := make([]*AttendanceResponse, len(attendances))
	for i, a := range attendances {
		resps[i] = ToAttendanceResponseFromDomain(&a)
	}
	return resps
}
