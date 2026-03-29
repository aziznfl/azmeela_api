package utils

import "math"

type PaginationMeta struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
}

func CalculatePagination(total int64, page, limit int) PaginationMeta {
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage == 0 {
		lastPage = 1
	}

	return PaginationMeta{
		Total:       total,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}
}

func GetOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * limit
}
