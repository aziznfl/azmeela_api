package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Fetch(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Transaction{})

	for k, v := range filter {
		if k == "search" {
			searchVal := "%" + v.(string) + "%"
			query = query.Joins("LEFT JOIN customers ON customers.id = transactions.customer_id").
				Where("LOWER(transactions.transaction_code) LIKE LOWER(?) OR LOWER(customers.full_name) LIKE LOWER(?)", searchVal, searchVal)
		} else {
			query = query.Where("transactions."+k+" = ?", v)
		}
	}

	query = query.Order("transactions.created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply Pagination & Joins (Senior Optimization: one query for 1:1 relations)
	err = query.
		Joins("Customer").
		Joins("Customer.CustomerType").
		Joins("Status").
		Offset(offset).Limit(limit).Order("transactions.id DESC").
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.WithContext(ctx).
		Joins("Customer").
		Joins("Customer.CustomerType").
		Joins("Status").
		Preload("Details.ProductPrice.Product.ProductCode.Type").
		Preload("Details.ProductPrice.Product.Variants").
		Preload("Details.ProductPrice.Size").
		First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Store(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *transactionRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, id).Error
}

func (r *transactionRepository) FetchStatuses(ctx context.Context) ([]domain.TransactionStatus, error) {
	var statuses []domain.TransactionStatus
	err := r.db.WithContext(ctx).Order("id ASC").Find(&statuses).Error
	return statuses, err
}
func (r *transactionRepository) FetchLogs(ctx context.Context, id int) ([]domain.TransactionLog, error) {
	var logs []domain.TransactionLog
	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", id).
		Joins("Admin").
		Joins("OldStatus").
		Joins("NewStatus").
		Order("transaction_logs.created_at DESC").
		Find(&logs).Error
	return logs, err
}
func (r *transactionRepository) GetLastTransactionCode(ctx context.Context, prefix string) (string, error) {
	var lastCode string
	err := r.db.WithContext(ctx).Table("transactions").
		Where("transaction_code LIKE ?", prefix+"%").
		Order("transaction_code DESC").
		Select("transaction_code").
		Limit(1).
		Row().Scan(&lastCode)

	// If no record found, it might return an error or empty string
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}
	return lastCode, nil
}
