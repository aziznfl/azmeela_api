package domain

import (
	"context"
	"time"
)

type ProductType struct {
	ID        int       `gorm:"primaryKey;column:id;autoIncrement"`
	Name      string    `gorm:"column:name;unique"`
	WebStatus int       `gorm:"column:web_status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ProductType) TableName() string {
	return "product_types"
}

type ProductSize struct {
	ID   int    `gorm:"primaryKey;column:id;autoIncrement"`
	Name string `gorm:"column:name;unique"`
}

func (ProductSize) TableName() string {
	return "product_sizes"
}

type ProductCode struct {
	ID            int            `gorm:"primaryKey;column:id;autoIncrement"`
	ProductTypeID int            `gorm:"column:product_type_id;uniqueIndex:idx_product_code_name"`
	Name          string         `gorm:"column:name;uniqueIndex:idx_product_code_name"`
	WebStatus     int            `gorm:"column:web_status"`
	CodeStatus    int            `gorm:"column:code_status"`
	Description   string         `gorm:"column:description"`
	Information   string         `gorm:"column:information"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	Type          *ProductType   `gorm:"foreignKey:ProductTypeID"`
	Products      []Product      `gorm:"foreignKey:ProductCodeID"`
}

func (ProductCode) TableName() string {
	return "product_codes"
}

type Product struct {
	ID            int            `gorm:"primaryKey;column:id;autoIncrement"`
	ProductCodeID int            `gorm:"column:product_code_id;uniqueIndex:idx_product_variation"`
	AdminID       int            `gorm:"column:admin_id"`
	Status        int            `gorm:"column:status"`
	SKU           string         `gorm:"column:sku"`
	Color         string         `gorm:"column:color;uniqueIndex:idx_product_variation"`
	Tags          string         `gorm:"column:tags"`
	WebStatus     int            `gorm:"column:web_status"`
	SEOLink       string         `gorm:"column:seo_link"`
	Views         int            `gorm:"column:views"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	ProductCode   *ProductCode   `gorm:"foreignKey:ProductCodeID"`
	Variants      []ProductPrice `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}

type ProductPrice struct {
	ID             int          `gorm:"primaryKey;column:id;autoIncrement"`
	ProductID      int          `gorm:"column:product_id;uniqueIndex:idx_product_price_unique"`
	CustomerTypeID int          `gorm:"column:customer_type_id;uniqueIndex:idx_product_price_unique"`
	ProductSizeID  int          `gorm:"column:product_size_id;uniqueIndex:idx_product_price_unique"`
	AdminID        int          `gorm:"column:admin_id"`
	Price          float64      `gorm:"column:price"`
	Specification  string       `gorm:"column:specification"`
	Stock          int          `gorm:"column:stock"`
	CartedCount    int          `gorm:"column:carted_count"`
	SoldCount      int          `gorm:"column:sold_count"`
	Weight         int          `gorm:"column:weight"`
	ProductDiscount int         `gorm:"column:product_discount"`
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	Product        *Product      `gorm:"foreignKey:ProductID"`
	Size           *ProductSize  `gorm:"foreignKey:ProductSizeID"`
	CustomerType   *CustomerType `gorm:"foreignKey:CustomerTypeID"`
}

type ProductStockLog struct {
	ID             int       `gorm:"primaryKey;column:id;autoIncrement"`
	ProductPriceID int       `gorm:"column:product_price_id"`
	Quantity       int       `gorm:"column:quantity"`
	AdminID        int       `gorm:"column:admin_id"`
	InputType      int       `gorm:"column:input_type;default:1"`
	CreatedAt      time.Time `gorm:"column:input_date;autoCreateTime"`
}

func (ProductStockLog) TableName() string {
	return "product_stock_inputs"
}

func (ProductPrice) TableName() string {
	return "product_prices"
}

type ProductRepository interface {
	FetchTypes(ctx context.Context) ([]ProductType, error)
	GetTypeByID(ctx context.Context, id int) (*ProductType, error)
	CreateType(ctx context.Context, pType *ProductType) error
	UpdateType(ctx context.Context, pType *ProductType) error
	DeleteType(ctx context.Context, id int) error

	FetchSizes(ctx context.Context) ([]ProductSize, error)
	
	FetchCodes(ctx context.Context, filter map[string]interface{}) ([]ProductCode, error)
	GetCodeByID(ctx context.Context, id int) (*ProductCode, error)
	CreateCode(ctx context.Context, code *ProductCode) error
	UpdateCode(ctx context.Context, code *ProductCode) error
	DeleteCode(ctx context.Context, id int) error

	FetchPrices(ctx context.Context, filter map[string]interface{}) ([]ProductPrice, error)
	UpdateStock(ctx context.Context, id int, quantity int) error
	GetStockLogs(ctx context.Context, productPriceID int) ([]ProductStockLog, error)

	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type ProductUsecase interface {
	GetInventoryList(ctx context.Context, filter map[string]interface{}) ([]ProductCode, error)
	
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
	UpdateStock(ctx context.Context, id int, quantity int) error
	GetStockLogs(ctx context.Context, productPriceID int) ([]ProductStockLog, error)
}
