package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type overtimeUsecase struct {
	overtimeRepo domain.OvertimeRepository
}

func NewOvertimeUsecase(o domain.OvertimeRepository) domain.OvertimeUsecase {
	return &overtimeUsecase{
		overtimeRepo: o,
	}
}

func (u *overtimeUsecase) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.OvertimeResponse, error) {
	overtimes, err := u.overtimeRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response []domain.OvertimeResponse
	for _, ot := range overtimes {
		response = append(response, domain.OvertimeResponse{
			ID:           ot.ID,
			EmployeeID:   ot.EmployeeID,
			EmployeeName: ot.EmployeeName,
			Date:         ot.Date,
			TimeIn:       ot.TimeIn,
			TimeOut:      ot.TimeOut,
			Status:       ot.Status,
			Description:  ot.Description,
		})
	}
	return response, nil
}

func (u *overtimeUsecase) RequestOvertime(ctx context.Context, employeeID int, req *domain.OvertimeRequest) (*domain.Overtime, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	ot := &domain.Overtime{
		EmployeeID:  employeeID,
		Date:        parsedDate,
		TimeIn:      req.TimeIn,
		TimeOut:     req.TimeOut,
		Status:      0, // Created
		Description: req.Description,
	}

	err = u.overtimeRepo.Store(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, nil
}

func (u *overtimeUsecase) UpdateStatus(ctx context.Context, id int, req *domain.OvertimeStatusUpdate) error {
	_, err := u.overtimeRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("overtime request not found")
	}

	return u.overtimeRepo.UpdateStatus(ctx, id, req.Status)
}
