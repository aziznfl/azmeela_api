package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type attendanceUsecase struct {
	attendanceRepo domain.AttendanceRepository
}

// NewAttendanceUsecase will create a new representation of domain.AttendanceUsecase interface
func NewAttendanceUsecase(a domain.AttendanceRepository) domain.AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: a,
	}
}

func mapToResponses(attendances []domain.Attendance) []domain.AttendanceResponse {
	var response []domain.AttendanceResponse
	for _, att := range attendances {
		response = append(response, domain.AttendanceResponse{
			ID:           att.ID,
			EmployeeID:   att.EmployeeID,
			EmployeeName: att.EmployeeName,
			Date:         att.Date,
			TimeIn:       att.TimeIn,
			TimeOut:      att.TimeOut,
			Location:     att.Location,
			Note:         att.Note,
		})
	}
	return response
}

func (u *attendanceUsecase) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.AttendanceResponse, error) {
	attendances, err := u.attendanceRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}
	return mapToResponses(attendances), nil
}

func (u *attendanceUsecase) GetTodayAttendances(ctx context.Context) ([]domain.AttendanceResponse, error) {
	attendances, err := u.attendanceRepo.GetTodayAttendances(ctx)
	if err != nil {
		return nil, err
	}
	return mapToResponses(attendances), nil
}

func (u *attendanceUsecase) ClockIn(ctx context.Context, employeeID int, req *domain.AttendanceRequest) (*domain.Attendance, error) {
	today := time.Now()
	nowStr := today.Format("15:04:05")

	// Check if already exist
	att, err := u.attendanceRepo.GetByDateAndEmployee(ctx, today, employeeID)
	if err == nil && att != nil {
		if att.TimeIn != "" {
			return nil, errors.New("already clocked in today")
		}

		att.TimeIn = nowStr
		if req != nil && req.Location != "" {
			att.Location = &req.Location
		}
		if req != nil && req.Note != "" {
			att.Note = &req.Note
		}
		err = u.attendanceRepo.Update(ctx, att)
		if err != nil {
			return nil, err
		}
		return att, nil
	}

	// Create new record for today
	newAtt := &domain.Attendance{
		EmployeeID: employeeID,
		Date:       today,
		TimeIn:     nowStr,
	}

	if req != nil && req.Location != "" {
		newAtt.Location = &req.Location
	}
	if req != nil && req.Note != "" {
		newAtt.Note = &req.Note
	}

	err = u.attendanceRepo.Store(ctx, newAtt)
	if err != nil {
		return nil, err
	}

	return newAtt, nil
}

func (u *attendanceUsecase) ClockOut(ctx context.Context, employeeID int, req *domain.AttendanceRequest) (*domain.Attendance, error) {
	today := time.Now()
	nowStr := today.Format("15:04:05")

	att, err := u.attendanceRepo.GetByDateAndEmployee(ctx, today, employeeID)
	if err != nil || att == nil {
		return nil, errors.New("must clock in first before checking out")
	}

	if att.TimeOut != nil {
		return nil, errors.New("already clocked out today")
	}

	att.TimeOut = &nowStr
	if req != nil && req.Location != "" {
		att.Location = &req.Location
	}
	if req != nil && req.Note != "" {
		att.Note = &req.Note
	}

	err = u.attendanceRepo.Update(ctx, att)
	if err != nil {
		return nil, err
	}

	return att, nil
}
