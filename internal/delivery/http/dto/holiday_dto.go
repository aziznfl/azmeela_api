package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type HolidayResponse struct {
	ID          int       `json:"id"`
	HolidayDate time.Time `json:"holiday_date"`
	Description string    `json:"description"`
	IsRecurring bool      `json:"is_recurring"`
}

type HolidayRequest struct {
	HolidayDate string `json:"holiday_date" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsRecurring bool   `json:"is_recurring"`
}

func ToHolidayResponse(h *domain.Holiday) *HolidayResponse {
	if h == nil {
		return nil
	}
	return &HolidayResponse{
		ID:          h.ID,
		HolidayDate: h.HolidayDate,
		Description: h.Description,
		IsRecurring: h.IsRecurring,
	}
}

func ToHolidayListResponse(items []domain.Holiday) []*HolidayResponse {
	resps := make([]*HolidayResponse, len(items))
	for i, item := range items {
		resps[i] = ToHolidayResponse(&item)
	}
	return resps
}
