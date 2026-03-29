package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type holidayRepository struct {
	BaseRepository[domain.Holiday]
}

func NewHolidayRepository(db *gorm.DB) domain.HolidayRepository {
	return &holidayRepository{
		BaseRepository: BaseRepository[domain.Holiday]{db: db},
	}
}

func (r *holidayRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Holiday, error) {
	// For holiday, we dont use pagination yet so offset=0 and limit=1000
	res, _, err := r.BaseRepository.Fetch(ctx, filter, 0, 1000)
	return res, err
}
// Delete, Update, Store and GetByID are already handled by BaseRepository if they match the signature.
// Wait. domain.HolidayRepository methods:
// Store(ctx context.Context, holiday *Holiday) error
// Update(ctx context.Context, holiday *Holiday) error
// GetByID(ctx context.Context, id int) (*Holiday, error)
// Delete(ctx context.Context, id int) error
// These match the BaseRepository signatures.
