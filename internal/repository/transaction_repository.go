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
			query = query.Joins("LEFT JOIN t_customer ON t_customer.id_customer = t_transaksi.id_customer").
				Where("LOWER(t_transaksi.kode_transaksi) LIKE LOWER(?) OR LOWER(t_customer.name_customer) LIKE LOWER(?)", searchVal, searchVal)
		} else if k == "status_id" {
			query = query.Where("t_transaksi.transaksi_status = ?", v)
		} else if k == "customer_id" {
			query = query.Where("t_transaksi.id_customer = ?", v)
		} else {
			query = query.Where("t_transaksi."+k+" = ?", v)
		}
	}

	query = query.Order("t_transaksi.tgl_transaksi DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply Pagination & Joins (Senior Optimization: one query for 1:1 relations)
	err = query.
		Joins("Customer").
		Joins("Customer.CustomerType").
		Joins("Status").
		Offset(offset).Limit(limit).Order("t_transaksi.id_transaksi DESC").
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

	// Populate virtual fields
	for i := range transaction.Details {
		if transaction.Details[i].ProductPrice != nil && transaction.Details[i].ProductPrice.Product != nil {
			transaction.Details[i].ProductID = transaction.Details[i].ProductPrice.ProductID
			transaction.Details[i].ProductCodeID = transaction.Details[i].ProductPrice.Product.ProductCodeID
		}
	}

	return &transaction, nil
}

func (r *transactionRepository) Store(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		// 1. Manual wipe of existing details
		if err := db.Where("id_transaksi = ?", tx.ID).Delete(&domain.TransactionDetail{}).Error; err != nil {
			return err
		}

		// 2. Manual insert of new details
		if len(tx.Details) > 0 {
			if err := db.Create(&tx.Details).Error; err != nil {
				return err
			}
		}

		// 3. Save the transaction main record, omitting "Details" to prevent GORM
		// from trying to "link/unlink" associations which triggers the NOT NULL error.
		return db.Omit("Details").Save(tx).Error
	})
}

func (r *transactionRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, id).Error
}

func (r *transactionRepository) FetchStatuses(ctx context.Context) ([]domain.TransactionStatus, error) {
	var statuses []domain.TransactionStatus
	err := r.db.WithContext(ctx).Order("id_transaksi_status ASC").Find(&statuses).Error
	return statuses, err
}
func (r *transactionRepository) FetchLogs(ctx context.Context, id int) ([]domain.TransactionLog, error) {
	var logs []domain.TransactionLog
	err := r.db.WithContext(ctx).
		Where("id_transaksi = ?", id).
		Joins("Admin").
		Joins("OldStatus").
		Joins("NewStatus").
		Order("t_transaksi_log.tgl_transkasi_log DESC").
		Find(&logs).Error
	return logs, err
}
func (r *transactionRepository) GetLastTransactionCode(ctx context.Context, prefix string) (string, error) {
	var lastCode string
	err := r.db.WithContext(ctx).Table("t_transaksi").
		Where("kode_transaksi LIKE ?", prefix+"%").
		Order("kode_transaksi DESC").
		Select("kode_transaksi").
		Limit(1).
		Row().Scan(&lastCode)

	// If no record found, it might return an error or empty string
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}
	return lastCode, nil
}
