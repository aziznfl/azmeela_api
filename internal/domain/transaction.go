package domain

import (
	"context"
	"time"
)

type TransactionStatus struct {
	ID        int    `gorm:"primaryKey;autoIncrement:false"`
	Name      string `gorm:"type:varchar(20);unique"`
	Privilege int
}

func (TransactionStatus) TableName() string {
	return "transaction_statuses"
}

type Transaction struct {
	ID              int                `gorm:"primaryKey;autoIncrement"`
	CustomerID      int
	Customer        *Customer          `gorm:"foreignKey:CustomerID"`
	AdminID         int
	Admin           *Employee          `gorm:"foreignKey:AdminID"`
	TransferDate    *time.Time         `gorm:"type:date;default:null"`
	ShippingDate    *time.Time         `gorm:"type:date;default:null"`
	StatusID        int                `gorm:"default:1"`
	Status          *TransactionStatus `gorm:"foreignKey:StatusID"`
	ShippingCost    int
	TrackingNumber  *string            `gorm:"type:varchar(50);default:null"`
	Courier         *string            `gorm:"type:varchar(30);default:null"`
	TransactionCode string             `gorm:"type:varchar(15);unique"`
	Total           int
	Address         string             `gorm:"type:text"`
	PaymentCode     int
	Discount        int                `gorm:"default:0"`
	DiscountNote    *string            `gorm:"type:varchar(200);default:null"`
	DiscountType    int                `gorm:"default:1"`
	IsReminded      bool               `gorm:"default:false"`
	CreatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID"`
	Logs            []TransactionLog    `gorm:"foreignKey:TransactionID"`
}

type TransactionDetail struct {
	ID             int           `gorm:"primaryKey;autoIncrement"`
	TransactionID  int
	ProductPriceID int
	ProductPrice   *ProductPrice `gorm:"foreignKey:ProductPriceID"`
	Quantity       int
	Price          int
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}

type TransactionLog struct {
	ID            int                `gorm:"primaryKey;autoIncrement"`
	TransactionID int
	AdminID       int
	Admin         *Employee          `gorm:"foreignKey:AdminID"`
	OldStatusID   int
	NewStatusID   int
	OldStatus     *TransactionStatus `gorm:"foreignKey:OldStatusID"`
	NewStatus     *TransactionStatus `gorm:"foreignKey:NewStatusID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (TransactionLog) TableName() string {
	return "transaction_logs"
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionRequest struct {
	CustomerID      int
	AdminID         int
	TransferDate    *time.Time
	ShippingDate    *time.Time
	StatusID        int
	ShippingCost    int
	TrackingNumber  *string
	Courier         *string
	TransactionCode string
	Total           int
	Address         string
	PaymentCode     int
	Discount        int
	DiscountNote    *string
	DiscountType    int
	IsReminded      bool
	Details         []TransactionDetailRequest
}

type TransactionDetailRequest struct {
	ProductPriceID int
	Quantity       int
	Price          int
}

type TransactionRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]Transaction, int64, error)
	GetByID(ctx context.Context, id int) (*Transaction, error)
	Store(ctx context.Context, tx *Transaction) error
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id int) error
	FetchStatuses(ctx context.Context) ([]TransactionStatus, error)
	FetchLogs(ctx context.Context, id int) ([]TransactionLog, error)
	GetLastTransactionCode(ctx context.Context, prefix string) (string, error)
}

type TransactionUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}, page, limit int) ([]Transaction, PaginationMeta, error)
	GetByID(ctx context.Context, id int) (*Transaction, error)
	Create(ctx context.Context, req *TransactionRequest) (*Transaction, error)
	Update(ctx context.Context, id int, req *TransactionRequest) (*Transaction, error)
	Delete(ctx context.Context, id int) error
	GetStatuses(ctx context.Context) ([]TransactionStatus, error)
	GetLogs(ctx context.Context, id int) ([]TransactionLog, error)
	GenerateTransactionCode(ctx context.Context) (string, error)
}
