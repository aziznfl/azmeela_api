package domain

import (
	"context"
	"time"
)

type TransactionStatus struct {
	ID        int    `gorm:"primaryKey;column:id_transaksi_status;autoIncrement"`
	Name      string `gorm:"column:nama_transaksi_status;type:varchar(20);unique"`
	Privilege int    `gorm:"column:privilege"`
}

func (TransactionStatus) TableName() string {
	return "t_transaksi_status"
}

type Transaction struct {
	ID              int                `gorm:"primaryKey;column:id_transaksi;autoIncrement"`
	CustomerID      int                `gorm:"column:id_customer"`
	Customer        *Customer          `gorm:"foreignKey:CustomerID"`
	AdminID         int                `gorm:"column:id_admin"`
	Admin           *Employee          `gorm:"foreignKey:AdminID"`
	TransactionDate time.Time          `gorm:"column:tgl_transaksi;type:timestamp"`
	CreatedAt       time.Time          `gorm:"column:tgl_cur_time;default:CURRENT_TIMESTAMP"`
	TransferDate    *time.Time         `gorm:"column:tgl_transfer;type:date;default:null"`
	ShippingDate    *time.Time         `gorm:"column:tgl_pengiriman;type:date;default:null"`
	StatusID        int                `gorm:"column:transaksi_status;default:1"`
	Status          *TransactionStatus `gorm:"foreignKey:StatusID"`
	ShippingCost    int                `gorm:"column:ongkir"`
	TrackingNumber  *string            `gorm:"column:resi;type:varchar(50);default:null"`
	Courier         *string            `gorm:"column:pengantar;type:varchar(30);default:null"`
	TransactionCode string             `gorm:"column:kode_transaksi;type:varchar(15);unique"`
	Total           int                `gorm:"column:total"`
	Address         string             `gorm:"column:alamat;type:text"`
	PaymentCode     int                `gorm:"column:kode_pembayaran"`
	Discount        int                `gorm:"column:diskon;default:0"`
	DiscountNote    *string            `gorm:"column:ket_diskon;type:varchar(200);default:null"`
	DiscountType    int                `gorm:"column:tipe_diskon;default:1"`
	IsReminded      int                `gorm:"column:reminded;default:0"`
	Details         []TransactionDetail `gorm:"foreignKey:TransactionID"`
	Logs            []TransactionLog    `gorm:"foreignKey:TransactionID"`
}

type TransactionDetail struct {
	ID             int           `gorm:"primaryKey;column:id_transaksi_detail;autoIncrement"`
	TransactionID  int           `gorm:"column:id_transaksi"`
	ProductPriceID int           `gorm:"column:id_barang"`
	ProductPrice   *ProductPrice `gorm:"foreignKey:ProductPriceID"`
	Quantity       int           `gorm:"column:qty"`
	Price          int           `gorm:"column:harga"`
	// Virtual fields for frontend convenience
	ProductCodeID int `gorm:"-"`
	ProductID     int `gorm:"-"`
}

func (TransactionDetail) TableName() string {
	return "t_transaksi_detail"
}

type TransactionLog struct {
	ID            int                `gorm:"primaryKey;column:id_transaksi_log;autoIncrement"`
	TransactionID int                `gorm:"column:id_transaksi"`
	AdminID       int                `gorm:"column:id_admin"`
	Admin         *Employee          `gorm:"foreignKey:AdminID"`
	OldStatusID   int                `gorm:"column:transaksi_status_old"`
	NewStatusID   int                `gorm:"column:transaksi_status_new"`
	OldStatus     *TransactionStatus `gorm:"foreignKey:OldStatusID"`
	NewStatus     *TransactionStatus `gorm:"foreignKey:NewStatusID"`
	CreatedAt     time.Time          `gorm:"column:tgl_transkasi_log;autoCreateTime"`
}

func (TransactionLog) TableName() string {
	return "t_transaksi_log"
}

func (Transaction) TableName() string {
	return "t_transaksi"
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
	IsReminded      int
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
