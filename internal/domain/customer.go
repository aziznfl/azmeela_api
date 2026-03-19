package domain

import (
	"context"
	"time"
)

type CustomerType struct {
	ID      int    `gorm:"primaryKey;column:id;autoIncrement"`
	Name    string `gorm:"column:name;type:varchar(100);not null;unique"`
	Initial string `gorm:"column:initial;type:varchar(1);not null;unique"`
}

func (CustomerType) TableName() string {
	return "customer_types"
}

type Customer struct {
	ID               int          `gorm:"primaryKey;column:id;autoIncrement"`
	CustomerTypeID   int          `gorm:"column:customer_type_id;not null;default:1"`
	CustomerType     CustomerType `gorm:"foreignKey:CustomerTypeID"`
	CreatedAt        time.Time    `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	FullName         string       `gorm:"column:full_name;type:varchar(100);not null"`
	PhoneNumber      *string      `gorm:"column:phone_number;type:varchar(16)"`
	Deposit          float64      `gorm:"column:deposit;not null;default:0"`
	AdminID          int          `gorm:"column:admin_id;not null;default:1"`
	MembershipStatus string       `gorm:"column:membership_status;type:membership_status_type;not null;default:'1'"`
	Username         *string      `gorm:"column:username;type:varchar(100);unique"`
	Email            *string      `gorm:"column:email;type:varchar(100)"`
	Password         *string      `gorm:"column:password;type:varchar(255)"`
	PasswordHash     *string      `gorm:"column:password_hash;type:varchar(255)"`
	IsActive         bool         `gorm:"column:is_active;not null;default:true;->:false"`
	Addresses        []CustomerAddress `gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return "customers"
}

type CustomerAddress struct {
	ID            int       `gorm:"primaryKey;column:id;autoIncrement"`
	CustomerID    int       `gorm:"column:customer_id;not null"`
	AdminID       *int      `gorm:"column:admin_id"`
	Country       string    `gorm:"column:country;type:varchar(100);not null"`
	Province      string    `gorm:"column:province;type:varchar(100);not null"`
	City          string    `gorm:"column:city;type:varchar(100);not null"`
	District      string    `gorm:"column:district;type:varchar(100);not null"`
	SubDistrict   *string   `gorm:"column:sub_district;type:varchar(100)"`
	StreetAddress string    `gorm:"column:street_address;type:varchar(200);not null"`
	PostalCode    *string   `gorm:"column:postal_code;type:varchar(11)"`
	SicepatID     *int      `gorm:"column:sicepat_id"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

func (CustomerAddress) TableName() string {
	return "customer_addresses"
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
