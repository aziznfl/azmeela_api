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
	ID              int                `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID      int                `json:"customer_id"`
	Customer        *Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	AdminID         int                `json:"admin_id"`
	Admin           *Employee          `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	TransferDate    *time.Time         `gorm:"type:date;default:null" json:"transfer_date"`
	ShippingDate    *time.Time         `gorm:"type:date;default:null" json:"shipping_date"`
	StatusID        int                `gorm:"default:1" json:"status_id"`
	Status          *TransactionStatus `gorm:"foreignKey:StatusID" json:"status,omitempty"`
	ShippingCost    int                `json:"shipping_cost"`
	TrackingNumber  *string            `gorm:"type:varchar(50);default:null" json:"tracking_number"`
	Courier         *string            `gorm:"type:varchar(30);default:null" json:"courier"`
	TransactionCode string             `gorm:"type:varchar(15);unique" json:"transaction_code"`
	Total           int                `json:"total"`
	Address         string             `gorm:"type:text" json:"address"`
	PaymentCode     int                `json:"payment_code"`
	Discount        int                `gorm:"default:0" json:"discount"`
	DiscountNote    *string            `gorm:"type:varchar(200);default:null" json:"discount_note"`
	DiscountType    int                `gorm:"default:1" json:"discount_type"`
	IsReminded      bool               `gorm:"default:false" json:"is_reminded"`
	CreatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
	Logs            []TransactionLog    `gorm:"foreignKey:TransactionID" json:"logs,omitempty"`
}

type TransactionDetail struct {
	ID             int           `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID  int           `json:"transaction_id"`
	ProductPriceID int           `json:"product_price_id"`
	ProductPrice   *ProductPrice `gorm:"foreignKey:ProductPriceID" json:"product_price,omitempty"`
	Quantity       int           `json:"quantity"`
	Price          int           `json:"price"`
	// Virtual fields for frontend convenience
	ProductCodeID int `gorm:"-" json:"product_code_id"`
	ProductID     int `gorm:"-" json:"product_id"`
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
