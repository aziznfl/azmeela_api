package domain

type PaginationMeta struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
}

func CalculatePagination(total int64, page, limit int) (int, PaginationMeta) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	return offset, PaginationMeta{
		Total:       total,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}
}
