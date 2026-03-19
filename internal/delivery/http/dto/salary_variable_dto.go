package dto

import (
	"github.com/azmeela/sispeg-api/internal/domain"
)

type SalaryVariableResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Value int    `json:"value"`
}

type SalaryVariableRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  int    `json:"type" binding:"required"`
	Value int    `json:"value" binding:"required,min=0"`
}

func ToSalaryVariableResponse(sv *domain.SalaryVariable) *SalaryVariableResponse {
	if sv == nil {
		return nil
	}
	return &SalaryVariableResponse{
		ID:    sv.ID,
		Name:  sv.Name,
		Type:  sv.Type,
		Value: sv.Value,
	}
}

func ToSalaryVariableListResponse(items []domain.SalaryVariable) []*SalaryVariableResponse {
	resps := make([]*SalaryVariableResponse, len(items))
	for i, item := range items {
		resps[i] = ToSalaryVariableResponse(&item)
	}
	return resps
}
