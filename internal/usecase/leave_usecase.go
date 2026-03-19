package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type leaveUsecase struct {
	leaveRepo domain.LeaveRepository
}

func NewLeaveUsecase(l domain.LeaveRepository) domain.LeaveUsecase {
	return &leaveUsecase{
		leaveRepo: l,
	}
}

func (u *leaveUsecase) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.LeaveResponse, error) {
	leaves, err := u.leaveRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response []domain.LeaveResponse
	for _, l := range leaves {
		response = append(response, domain.LeaveResponse{
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
		})
	}
	return response, nil
}

func (u *leaveUsecase) RequestLeave(ctx context.Context, employeeID int, req *domain.LeaveRequest) (*domain.Leave, error) {
	parsedDate, err := time.Parse("2006-01-02", req.LeaveDate)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	leave := &domain.Leave{
		EmployeeID:  employeeID,
		Type:        req.Type,
		LeaveDate:   parsedDate,
		Durations:   req.Durations,
		Status:      0, // created
		Description: req.Description,
	}

	err = u.leaveRepo.Store(ctx, leave)
	if err != nil {
		return nil, err
	}

	return leave, nil
}

func (u *leaveUsecase) UpdateStatus(ctx context.Context, id int, req *domain.LeaveStatusUpdate) error {
	_, err := u.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("leave request not found")
	}

	return u.leaveRepo.UpdateStatus(ctx, id, req.Status)
}
