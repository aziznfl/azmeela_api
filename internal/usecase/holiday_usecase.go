package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type holidayUsecase struct {
	holidayRepo domain.HolidayRepository
}

func NewHolidayUsecase(h domain.HolidayRepository) domain.HolidayUsecase {
	return &holidayUsecase{
		holidayRepo: h,
	}
}

func (u *holidayUsecase) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Holiday, error) {
	return u.holidayRepo.Fetch(ctx, filter)
}

func (u *holidayUsecase) Create(ctx context.Context, req *domain.HolidayRequest) (*domain.Holiday, error) {
	parsedDate, err := time.Parse("2006-01-02", req.HolidayDate)
	if err != nil {
		return nil, errors.New("invalid holiday_date format, expected YYYY-MM-DD")
	}

	hol := &domain.Holiday{
		HolidayDate: parsedDate,
		Month:       int(parsedDate.Month()),
		Day:         int(parsedDate.Day()),
		Description: req.Description,
		IsRecurring: req.IsRecurring,
	}

	err = u.holidayRepo.Store(ctx, hol)
	if err != nil {
		return nil, err
	}
	return hol, nil
}
func (u *holidayUsecase) Update(ctx context.Context, id int, req *domain.HolidayRequest) (*domain.Holiday, error) {
	hol, err := u.holidayRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse("2006-01-02", req.HolidayDate)
	if err != nil {
		return nil, errors.New("invalid holiday_date format, expected YYYY-MM-DD")
	}

	hol.HolidayDate = parsedDate
	hol.Month = int(parsedDate.Month())
	hol.Day = int(parsedDate.Day())
	hol.Description = req.Description
	hol.IsRecurring = req.IsRecurring

	err = u.holidayRepo.Update(ctx, hol)
	if err != nil {
		return nil, err
	}
	return hol, nil
}

func (u *holidayUsecase) Delete(ctx context.Context, id int) error {
	return u.holidayRepo.Delete(ctx, id)
}
