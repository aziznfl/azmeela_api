package domain

import (
	"context"
	"time"
)

type ProductType struct {
	ID        int       `gorm:"primaryKey;column:id_product_type;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name_product_type;unique" json:"name"`
	WebStatus int       `gorm:"column:status_web" json:"web_status"`
	CreatedAt time.Time `gorm:"column:date;autoCreateTime" json:"created_at"`
}

func (ProductType) TableName() string {
	return "t_product_type"
}

type ProductSize struct {
	ID   int    `gorm:"primaryKey;column:id_product_size;autoIncrement"`
	Name string `gorm:"column:name_product_size;unique"`
}

func (ProductSize) TableName() string {
	return "t_product_size"
}

type ProductCode struct {
	ID            int            `gorm:"primaryKey;column:id_product_code;autoIncrement" json:"id"`
	ProductTypeID int            `gorm:"column:id_product_type;uniqueIndex:idx_product_code_name" json:"product_type_id"`
	Name          string         `gorm:"column:name_product_code;uniqueIndex:idx_product_code_name" json:"name"`
	WebStatus     int            `gorm:"column:status_web" json:"web_status"`
	CodeStatus    int            `gorm:"column:status_code" json:"code_status"`
	Description   string         `gorm:"column:description" json:"description"`
	Information   string         `gorm:"column:information" json:"information"`
	CreatedAt     time.Time      `gorm:"column:date;autoCreateTime" json:"created_at"`
	Type          *ProductType   `gorm:"foreignKey:ProductTypeID" json:"type,omitempty"`
	Products      []Product      `gorm:"foreignKey:ProductCodeID" json:"products,omitempty"`
}

func (ProductCode) TableName() string {
	return "t_product_code"
}

type ProductCodeWithType struct {
	ID              int    `json:"id"`
	ProductTypeID   int    `json:"product_type_id"`
	Name            string `json:"name"`
	ProductTypeName string `json:"product_type_name"`
}

type Product struct {
	ID            int            `gorm:"primaryKey;column:id_product;autoIncrement" json:"id"`
	ProductCodeID int            `gorm:"column:id_product_code;uniqueIndex:idx_product_variation" json:"product_code_id"`
	AdminID       int            `gorm:"column:id_admin" json:"admin_id"`
	Status        int            `gorm:"column:status_product" json:"status"`
	SKU           string         `gorm:"column:sku" json:"sku"`
	Color         string         `gorm:"column:color;uniqueIndex:idx_product_variation" json:"color"`
	Tags          string         `gorm:"column:tags" json:"tags"`
	WebStatus     int            `gorm:"column:status_web" json:"web_status"`
	SEOLink       string         `gorm:"column:seo_link" json:"seo_link"`
	Views         int            `gorm:"column:viewed" json:"views"`
	CreatedAt     time.Time      `gorm:"column:date;autoCreateTime" json:"created_at"`
	ProductCode   *ProductCode   `gorm:"foreignKey:ProductCodeID" json:"product_code,omitempty"`
	Variants      []ProductPrice `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
}

func (Product) TableName() string {
	return "t_product"
}

type ProductPrice struct {
	ID             int          `gorm:"primaryKey;column:id_product_price;autoIncrement" json:"id"`
	ProductID      int          `gorm:"column:id_product;uniqueIndex:idx_product_price_unique" json:"product_id"`
	CustomerTypeID int          `gorm:"column:id_customer_type;uniqueIndex:idx_product_price_unique" json:"customer_type_id"`
	ProductSizeID  int          `gorm:"column:id_product_size;uniqueIndex:idx_product_price_unique" json:"product_size_id"`
	AdminID        int          `gorm:"column:id_admin" json:"admin_id"`
	Price          float64      `gorm:"column:price" json:"price"`
	Specification  string       `gorm:"column:spesification" json:"specification"`
	Stock          int          `gorm:"column:stock" json:"stock"`
	CartedCount    int          `gorm:"column:carted" json:"carted_count"`
	SoldCount      int          `gorm:"column:buyed" json:"sold_count"`
	Weight         int          `gorm:"column:weight" json:"weight"`
	ProductDiscount int         `gorm:"column:diskon_produk" json:"product_discount"`
	CreatedAt      time.Time    `gorm:"column:date;autoCreateTime" json:"created_at"`
	Product        *Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Size           *ProductSize  `gorm:"foreignKey:ProductSizeID" json:"size,omitempty"`
	CustomerType   *CustomerType `gorm:"foreignKey:CustomerTypeID" json:"customer_type,omitempty"`
}

type ProductStockLog struct {
	ID             int           `gorm:"primaryKey;column:id_product_input_stok;autoIncrement"`
	ProductPriceID int           `gorm:"column:id_product_price"`
	ProductID      int           `gorm:"-"`
	Quantity       int           `gorm:"column:qty"`
	AdminID        int           `gorm:"column:id_admin"`
	InputType      int       `gorm:"column:type_input;default:1"`
	CreatedAt      time.Time `gorm:"column:date_input_product;autoCreateTime"`
	Admin          *Employee     `gorm:"foreignKey:AdminID"`
	ProductPrice   *ProductPrice `gorm:"foreignKey:ProductPriceID"`
}

func (ProductStockLog) TableName() string {
	return "t_product_input_stock"
}

func (ProductPrice) TableName() string {
	return "t_product_price"
}

type ProductRepository interface {
	FetchTypes(ctx context.Context) ([]ProductType, error)
	GetTypeByID(ctx context.Context, id int) (*ProductType, error)
	CreateType(ctx context.Context, pType *ProductType) error
	UpdateType(ctx context.Context, pType *ProductType) error
	DeleteType(ctx context.Context, id int) error

	FetchSizes(ctx context.Context) ([]ProductSize, error)
	
	FetchCodes(ctx context.Context, filter map[string]interface{}) ([]ProductCode, error)
	FetchCodesWithTypes(ctx context.Context, filter map[string]interface{}) ([]ProductCodeWithType, error)
	GetCodeByID(ctx context.Context, id int) (*ProductCode, error)
	CreateCode(ctx context.Context, code *ProductCode) error
	UpdateCode(ctx context.Context, code *ProductCode) error
	DeleteCode(ctx context.Context, id int) error

	FetchPrices(ctx context.Context, filter map[string]interface{}) ([]ProductPrice, error)
	GetPriceByID(ctx context.Context, id int) (*ProductPrice, error)
	UpdateStock(ctx context.Context, id int, quantity int, adminID int) error
	GetStockLogs(ctx context.Context, productPriceID int) ([]ProductStockLog, error)

	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id int) error
	FetchColors(ctx context.Context, productCodeID int) ([]ProductColorResponse, error)
	FetchSizesType(ctx context.Context, productID int, customerTypeID int) ([]ProductSizeTypeResponse, error)
}

type ProductUsecase interface {
	GetInventoryList(ctx context.Context, filter map[string]interface{}) ([]ProductCode, error)
	GetCodesWithTypes(ctx context.Context, filter map[string]interface{}) ([]ProductCodeWithType, error)
	
	GetProductTypes(ctx context.Context) ([]ProductType, error)
	CreateProductType(ctx context.Context, req *ProductType) error
	UpdateProductType(ctx context.Context, req *ProductType) error
	DeleteProductType(ctx context.Context, id int) error

	CreateProductCode(ctx context.Context, req *ProductCode) error
	UpdateProductCode(ctx context.Context, req *ProductCode) error
	DeleteProductCode(ctx context.Context, id int) error

	CreateProduct(ctx context.Context, req *Product) error
	UpdateProduct(ctx context.Context, req *Product) error
	DeleteProduct(ctx context.Context, id int) error

	GetProductSizes(ctx context.Context) ([]ProductSize, error)
	UpdateStock(ctx context.Context, id int, quantity int, adminID int) error
	GetStockLogs(ctx context.Context, productPriceID int) ([]ProductStockLog, error)
	GetProductColors(ctx context.Context, productCodeID int) ([]ProductColorResponse, error)
	GetProductSizesType(ctx context.Context, productID int, customerTypeID int) ([]ProductSizeTypeResponse, error)
}
