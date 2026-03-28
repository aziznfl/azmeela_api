package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type transactionUsecase struct {
	txRepo      domain.TransactionRepository
	productRepo domain.ProductRepository
}

func NewTransactionUsecase(r domain.TransactionRepository, pr domain.ProductRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		txRepo:      r,
		productRepo: pr,
	}
}

func (u *transactionUsecase) Fetch(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.Transaction, domain.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	transactions, total, err := u.txRepo.Fetch(ctx, filter, offset, limit)
	if err != nil {
		return nil, domain.PaginationMeta{}, err
	}

	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	meta := domain.PaginationMeta{
		Total:       total,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}

	return transactions, meta, nil
}

func (u *transactionUsecase) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	return u.txRepo.GetByID(ctx, id)
}

func (u *transactionUsecase) Create(ctx context.Context, req *domain.TransactionRequest) (*domain.Transaction, error) {
	details := make([]domain.TransactionDetail, len(req.Details))
	for i, d := range req.Details {
		details[i] = domain.TransactionDetail{
			ProductPriceID: d.ProductPriceID,
			Quantity:       d.Quantity,
			Price:          d.Price,
		}
	}

	tx := &domain.Transaction{
		CustomerID:      req.CustomerID,
		AdminID:         req.AdminID,
		TransferDate:    req.TransferDate,
		ShippingDate:    req.ShippingDate,
		StatusID:        req.StatusID,
		ShippingCost:    req.ShippingCost,
		TrackingNumber:  req.TrackingNumber,
		Courier:         req.Courier,
		TransactionCode: req.TransactionCode,
		Total:           req.Total,
		Address:         req.Address,
		PaymentCode:     req.PaymentCode,
		Discount:        req.Discount,
		DiscountNote:    req.DiscountNote,
		DiscountType:    req.DiscountType,
		IsReminded:      req.IsReminded,
		Details:         details,
	}

	// Default status if empty
	if tx.StatusID == 0 {
		tx.StatusID = 1 // default to new/unpaid based on schema
	}

	// txRepo.Store will handle the creation. No initial log needed as per request.

	err := u.txRepo.Store(ctx, tx)
	if err != nil {
		return nil, err
	}

	// Deduct stock for each item in transaction
	for _, d := range tx.Details {
		_ = u.productRepo.UpdateStock(ctx, d.ProductPriceID, -d.Quantity)
	}

	return u.txRepo.GetByID(ctx, tx.ID)
}

func (u *transactionUsecase) Update(ctx context.Context, id int, req *domain.TransactionRequest) (*domain.Transaction, error) {
	tx, err := u.txRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Log status change if status is different and provided
	if req.StatusID > 0 && tx.StatusID != req.StatusID {
		log := domain.TransactionLog{
			TransactionID: tx.ID,
			AdminID:       req.AdminID,
			OldStatusID:   tx.StatusID,
			NewStatusID:   req.StatusID,
			CreatedAt:     time.Now(),
		}
		tx.Logs = append(tx.Logs, log)
	}

	// Immutability checks
	if tx.TransactionCode != req.TransactionCode {
		return nil, errors.New("kode transaksi tidak dapat diubah")
	}
	if tx.CustomerID != req.CustomerID {
		return nil, errors.New("pelanggan tidak dapat diubah")
	}

	// Status-based constraints
	if tx.StatusID != 1 { // 1 is "Belum Dibayar"
		// If not waiting payment, only allow courier info changes
		tx.TrackingNumber = req.TrackingNumber
		tx.Courier = req.Courier
		tx.ShippingCost = req.ShippingCost
		tx.ShippingDate = req.ShippingDate
		// Also allow status change if provided
		if req.StatusID > 0 {
			tx.StatusID = req.StatusID
		}
	} else {
		// Revert old stock before updating to new details
		for _, d := range tx.Details {
			_ = u.productRepo.UpdateStock(ctx, d.ProductPriceID, d.Quantity)
		}

		// Full update allowed if still "Belum Dibayar"
		tx.TransferDate = req.TransferDate
		tx.ShippingDate = req.ShippingDate
		if req.StatusID > 0 {
			tx.StatusID = req.StatusID
		}
		tx.ShippingCost = req.ShippingCost
		tx.TrackingNumber = req.TrackingNumber
		tx.Courier = req.Courier
		tx.Total = req.Total
		tx.Address = req.Address
		tx.PaymentCode = req.PaymentCode
		tx.Discount = req.Discount
		tx.DiscountNote = req.DiscountNote
		tx.DiscountType = req.DiscountType
		tx.IsReminded = req.IsReminded

		// Update Details
		newDetails := make([]domain.TransactionDetail, len(req.Details))
		for i, d := range req.Details {
			newDetails[i] = domain.TransactionDetail{
				TransactionID:  tx.ID,
				ProductPriceID: d.ProductPriceID,
				Quantity:       d.Quantity,
				Price:          d.Price,
			}
		}
		tx.Details = newDetails
	}

	// Senior Trick: Clear associations pointers before Save to prevent GORM 
	// from using stale loaded objects to overwrite our ID changes.
	tx.Status = nil
	tx.Customer = nil
	tx.Admin = nil

	err = u.txRepo.Update(ctx, tx)
	if err != nil {
		return nil, err
	}

	// If details were updated (Status was 1), deduct new stock
	if tx.StatusID == 1 {
		for _, d := range tx.Details {
			_ = u.productRepo.UpdateStock(ctx, d.ProductPriceID, -d.Quantity)
		}
	}

	return u.txRepo.GetByID(ctx, tx.ID)
}

func (u *transactionUsecase) Delete(ctx context.Context, id int) error {
	return u.txRepo.Delete(ctx, id)
}

func (u *transactionUsecase) GetStatuses(ctx context.Context) ([]domain.TransactionStatus, error) {
	return u.txRepo.FetchStatuses(ctx)
}

func (u *transactionUsecase) GetLogs(ctx context.Context, id int) ([]domain.TransactionLog, error) {
	return u.txRepo.FetchLogs(ctx, id)
}

func (u *transactionUsecase) GenerateTransactionCode(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("OD%d%02d", now.Year(), now.Month())
	
	lastCode, err := u.txRepo.GetLastTransactionCode(ctx, prefix)
	if err != nil {
		return "", err
	}

	increment := 1
	if lastCode != "" && len(lastCode) >= 10 {
		// Extract numeric part (last 4 digits)
		fmt.Sscanf(lastCode[len(lastCode)-4:], "%d", &increment)
		increment++
	}

	return fmt.Sprintf("%s%04d", prefix, increment), nil
}
