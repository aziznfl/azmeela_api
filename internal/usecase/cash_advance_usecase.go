package usecase

import (
	"context"
	"errors"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type cashAdvanceUsecase struct {
	caRepo domain.CashAdvanceRepository
}

func NewCashAdvanceUsecase(c domain.CashAdvanceRepository) domain.CashAdvanceUsecase {
	return &cashAdvanceUsecase{
		caRepo: c,
	}
}

func (u *cashAdvanceUsecase) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.CashAdvanceResponse, error) {
	results, err := u.caRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response []domain.CashAdvanceResponse
	for _, ca := range results {
		response = append(response, domain.CashAdvanceResponse{
			ID:           ca.ID,
			EmployeeID:   ca.EmployeeID,
			EmployeeName: ca.EmployeeName,
			Amount:       ca.Amount,
			Purpose:      ca.Purpose,
			Status:       ca.Status,
			CreatedAt:    ca.CreatedAt,
		})
	}
	return response, nil
}

func (u *cashAdvanceUsecase) RequestCashAdvance(ctx context.Context, employeeID int, req *domain.CashAdvanceRequest) (*domain.CashAdvance, error) {
	ca := &domain.CashAdvance{
		EmployeeID: employeeID,
		Amount:     req.Amount,
		Purpose:    req.Purpose,
		Status:     0, // Created
	}

	err := u.caRepo.Store(ctx, ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (u *cashAdvanceUsecase) UpdateStatus(ctx context.Context, id int, req *domain.CashAdvanceStatusUpdate) error {
	ca, err := u.caRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("cash advance request not found")
	}

	// Update status
	err = u.caRepo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return err
	}

	// If approved, store in history as debt
	if req.Status == 1 {
		history := &domain.CashAdvanceHistory{
			EmployeeID: ca.EmployeeID,
			Date:       ca.CreatedAt,
			Amount:     ca.Amount,
			Type:       1, // 1: utang / debt
		}

		// Attempt to store, but if it fails we might want to log it
		u.caRepo.StoreHistory(ctx, history)
	}

	return nil
}

func (u *cashAdvanceUsecase) AddPayment(ctx context.Context, req *domain.CashAdvancePayment) error {
	history := &domain.CashAdvanceHistory{
		EmployeeID: req.EmployeeID,
		Date:       req.Date,
		Amount:     req.Amount,
		Type:       2, // 2: bayar / payment
	}

	return u.caRepo.StoreHistory(ctx, history)
}
