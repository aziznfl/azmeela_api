package domain

import (
	"context"
	"time"
)

type CustomerType struct {
	ID      int    `gorm:"primaryKey;column:id_customer_type;autoIncrement"`
	Name    string `gorm:"column:name_customer_type;type:varchar(100);not null"`
	Initial string `gorm:"column:initial_customer_type;type:varchar(1);not null"`
}

func (CustomerType) TableName() string {
	return "t_customer_type"
}

type Customer struct {
	ID               int          `gorm:"primaryKey;column:id_customer;autoIncrement"`
	CustomerTypeID   int          `gorm:"column:id_customer_type;not null;default:1"`
	CustomerType     CustomerType `gorm:"foreignKey:CustomerTypeID"`
	CreatedAt        time.Time    `gorm:"column:date_sign;not null;default:CURRENT_TIMESTAMP"`
	FullName         string       `gorm:"column:name_customer;type:varchar(100);not null"`
	PhoneNumber      *string      `gorm:"column:no_telepon;type:varchar(16)"`
	Deposit          float64      `gorm:"column:deposito;not null;default:0"`
	AdminID          int          `gorm:"column:id_admin;not null;default:1"`
	MembershipStatus string       `gorm:"column:status_keanggotaan;type:text;not null;default:'1'"`
	Username         *string      `gorm:"column:username;type:varchar(100);unique"`
	Email            *string      `gorm:"column:email;type:varchar(100)"`
	Password         *string      `gorm:"column:password;type:varchar(255)"`
	PasswordHash     *string      `gorm:"column:password_hash;type:varchar(255)"`
	Addresses        []CustomerAddress `gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return "t_customer"
}

type CustomerAddress struct {
	ID            int       `gorm:"primaryKey;column:id_customer_alamat;autoIncrement"`
	CustomerID    int       `gorm:"column:id_customer;not null"`
	AdminID       *int      `gorm:"column:id_admin"`
	Country       string    `gorm:"column:negara;type:varchar(100);not null"`
	Province      string    `gorm:"column:provinsi;type:varchar(100);not null"`
	City          string    `gorm:"column:kota;type:varchar(100);not null"`
	District      string    `gorm:"column:kecamatan;type:varchar(100);not null"`
	SubDistrict   *string   `gorm:"column:kelurahan;type:varchar(100)"`
	StreetAddress string    `gorm:"column:jalan;type:varchar(200);not null"`
	PostalCode    *string   `gorm:"column:kode_pos;type:varchar(11)"`
	SicepatID     *int      `gorm:"column:id_sicepat"`
}

func (CustomerAddress) TableName() string {
	return "t_customer_alamat"
}

type CustomerAddressRequest struct {
	CustomerID    int     `json:"customer_id"`
	Country       string  `json:"country" binding:"required"`
	Province      string  `json:"province" binding:"required"`
	City          string  `json:"city" binding:"required"`
	District      string  `json:"district" binding:"required"`
	SubDistrict   *string `json:"sub_district"`
	StreetAddress string  `json:"street_address" binding:"required"`
	PostalCode    *string `json:"postal_code"`
	SicepatID     *int    `json:"sicepat_id"`
}

type CustomerRequest struct {
	CustomerTypeID   int     `json:"customer_type_id"`
	FullName         string  `json:"full_name" binding:"required"`
	PhoneNumber      *string `json:"phone_number"`
	Deposit          float64 `json:"deposit"`
	MembershipStatus string  `json:"membership_status"`
	Username         *string `json:"username"`
	Email            *string `json:"email"`
	Password         *string `json:"password"`
	IsActive         *bool   `json:"is_active"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
}

type CustomerRepository interface {
	Fetch(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]Customer, int64, error)
	GetByID(ctx context.Context, id int) (*Customer, error)
	GetByUsername(ctx context.Context, username string) (*Customer, error)
	Store(ctx context.Context, customer *Customer) error
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id int) error

	// Types
	FetchTypes(ctx context.Context) ([]CustomerType, error)
	GetTypeByID(ctx context.Context, id int) (*CustomerType, error)
	CreateType(ctx context.Context, cType *CustomerType) error
	UpdateType(ctx context.Context, cType *CustomerType) error
	DeleteType(ctx context.Context, id int) error

	// Addresses
	FetchAddresses(ctx context.Context, customerID int) ([]CustomerAddress, error)
	GetAddressByID(ctx context.Context, id int) (*CustomerAddress, error)
	StoreAddress(ctx context.Context, address *CustomerAddress) error
	UpdateAddress(ctx context.Context, address *CustomerAddress) error
	DeleteAddress(ctx context.Context, id int) error
}

type CustomerUsecase interface {
	Fetch(ctx context.Context, filter map[string]interface{}, page, limit int) ([]Customer, PaginationMeta, error)
	GetByID(ctx context.Context, id int) (*Customer, error)
	Create(ctx context.Context, req *CustomerRequest) (*Customer, error)
	Update(ctx context.Context, id int, req *CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id int) error

	GetTypes(ctx context.Context) ([]CustomerType, error)
	CreateType(ctx context.Context, req *CustomerType) error
	UpdateType(ctx context.Context, id int, req *CustomerType) error
	DeleteType(ctx context.Context, id int) error

	// Addresses
	GetAddresses(ctx context.Context, customerID int) ([]CustomerAddress, error)
	CreateAddress(ctx context.Context, req *CustomerAddressRequest) (*CustomerAddress, error)
	UpdateAddress(ctx context.Context, id int, req *CustomerAddressRequest) (*CustomerAddress, error)
	DeleteAddress(ctx context.Context, id int) error
}
