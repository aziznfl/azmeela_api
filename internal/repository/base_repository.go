package repository

import (
	"context"
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Fetch(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]T, int64, error) {
	var results []T
	var total int64
	query := r.db.WithContext(ctx).Model(new(T))

	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&results).Error
	return results, total, err
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	var result T
	err := r.db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) Store(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}
