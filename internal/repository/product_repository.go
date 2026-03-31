package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FetchTypes(ctx context.Context) ([]domain.ProductType, error) {
	var types []domain.ProductType
	err := r.db.WithContext(ctx).Order("id_product_type ASC").Find(&types).Error
	return types, err
}

func (r *productRepository) GetTypeByID(ctx context.Context, id int) (*domain.ProductType, error) {
	var pType domain.ProductType
	err := r.db.WithContext(ctx).First(&pType, id).Error
	return &pType, err
}

func (r *productRepository) CreateType(ctx context.Context, pType *domain.ProductType) error {
	return r.db.WithContext(ctx).Create(pType).Error
}

func (r *productRepository) UpdateType(ctx context.Context, pType *domain.ProductType) error {
	return r.db.WithContext(ctx).Save(pType).Error
}

func (r *productRepository) DeleteType(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductType{}, id).Error
}

func (r *productRepository) FetchSizes(ctx context.Context) ([]domain.ProductSize, error) {
	var sizes []domain.ProductSize
	err := r.db.WithContext(ctx).Order("id ASC").Find(&sizes).Error
	return sizes, err
}

func (r *productRepository) FetchCodes(ctx context.Context, filter map[string]interface{}) ([]domain.ProductCode, error) {
	var codes []domain.ProductCode

	query := r.db.WithContext(ctx).Model(&domain.ProductCode{})
	query = query.Preload("Type").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("color ASC")
		})

	if custTypeID, ok := filter["customer_type_id"].(int); ok && custTypeID != 0 {
		query = query.Where("EXISTS (SELECT 1 FROM t_product p JOIN t_product_price pp ON p.id_product = pp.id_product WHERE p.id_product_code = t_product_code.id_product_code AND pp.id_customer_type = ?)", custTypeID)
		query = query.Preload("Products.Variants", "id_customer_type = ?", custTypeID)
	} else {
		query = query.Preload("Products.Variants")
	}

	query = query.Preload("Products.Variants.Size")
	query = query.Preload("Products.Variants.CustomerType")

	if name, ok := filter["name"].(string); ok && name != "" {
		query = query.Where("name_product_code ILIKE ?", "%"+name+"%")
	}

	if typeID, ok := filter["product_type_id"].(int); ok && typeID != 0 {
		query = query.Where("id_product_type = ?", typeID)
	}

	if codeID, ok := filter["product_code_id"].(int); ok && codeID != 0 {
		query = query.Where("id_product_code = ?", codeID)
	}

	err := query.Order("date DESC").Find(&codes).Error
	return codes, err
}

func (r *productRepository) FetchCodesWithTypes(ctx context.Context, filter map[string]interface{}) ([]domain.ProductCodeWithType, error) {
	var results []domain.ProductCodeWithType
	query := r.db.WithContext(ctx).Table("t_product_code").
		Select("t_product_code.id_product_code as id, t_product_code.id_product_type as product_type_id, t_product_code.name_product_code as name, t_product_type.name_product_type as product_type_name").
		Joins("LEFT JOIN t_product_type ON t_product_code.id_product_type = t_product_type.id_product_type")

	if custTypeID, ok := filter["customer_type_id"].(int); ok && custTypeID != 0 {
		query = query.Where("EXISTS (SELECT 1 FROM t_product p JOIN t_product_price pp ON p.id_product = pp.id_product WHERE p.id_product_code = t_product_code.id_product_code AND pp.id_customer_type = ?)", custTypeID)
	}

	err := query.Order("LOWER(t_product_code.name_product_code) ASC").Scan(&results).Error
	return results, err
}

func (r *productRepository) GetCodeByID(ctx context.Context, id int) (*domain.ProductCode, error) {
	var code domain.ProductCode
	err := r.db.WithContext(ctx).
		Preload("Type").
		Preload("Products").
		Preload("Products.Variants").
		Preload("Products.Variants.Size").
		First(&code, id).Error
	return &code, err
}

func (r *productRepository) CreateCode(ctx context.Context, code *domain.ProductCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *productRepository) UpdateCode(ctx context.Context, code *domain.ProductCode) error {
	return r.db.WithContext(ctx).Model(code).Select(
		"Name", "ProductTypeID", "Description", "Information", "WebStatus", "CodeStatus",
	).Updates(code).Error
}

func (r *productRepository) DeleteCode(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductCode{}, id).Error
}

func (r *productRepository) FetchPrices(ctx context.Context, filter map[string]interface{}) ([]domain.ProductPrice, error) {
	var prices []domain.ProductPrice
	query := r.db.WithContext(ctx).Model(&domain.ProductPrice{}).
		Joins("Product").
		Joins("Product.ProductCode").
		Joins("Size")

	if productID, ok := filter["product_id"].(int); ok && productID != 0 {
		query = query.Where("t_product_price.id_product = ?", productID)
	}

	if customerTypeID, ok := filter["customer_type_id"].(int); ok && customerTypeID != 0 {
		query = query.Where("t_product_price.id_customer_type = ?", customerTypeID)
	}

	err := query.Order("t_product_price.date DESC").Find(&prices).Error
	return prices, err
}
func (r *productRepository) GetPriceByID(ctx context.Context, id int) (*domain.ProductPrice, error) {
	var pp domain.ProductPrice
	err := r.db.WithContext(ctx).First(&pp, id).Error
	return &pp, err
}

func (r *productRepository) UpdateStock(ctx context.Context, id int, quantity int, adminID int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current product_id from the product_price
		var pp domain.ProductPrice
		if err := tx.Select("id_product").First(&pp, id).Error; err != nil {
			return err
		}

		// 1. Update stock in product_prices
		if err := tx.Model(&domain.ProductPrice{}).Where("id_product_price = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
			return err
		}

		// 2. Insert into product_stock_inputs (ProductStockLog)
		log := domain.ProductStockLog{
			ProductPriceID: id,
			ProductID:      pp.ProductID,
			Quantity:       quantity,
			AdminID:        adminID,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *productRepository) GetStockLogs(ctx context.Context, productPriceID int) ([]domain.ProductStockLog, error) {
	var logs []domain.ProductStockLog
	err := r.db.WithContext(ctx).
		Joins("Admin").
		Joins("ProductPrice").
		Joins("ProductPrice.Size").
		Where("t_product_input_stock.id_product_price = ?", productPriceID).
		Order("t_product_input_stock.date_input_product DESC").Find(&logs).Error
	return logs, err
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Model(product).Select(
		"SKU", "Color", "Status", "WebStatus", "Tags", "SEOLink",
	).Updates(product).Error
}

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepository) FetchColors(ctx context.Context, productCodeID int) ([]domain.ProductColorResponse, error) {
	var results []domain.ProductColorResponse
	err := r.db.WithContext(ctx).Table("t_product").
		Select("id_product as id, id_product_code as product_code_id, color").
		Where("id_product_code = ?", productCodeID).
		Order("LOWER(color) ASC").
		Scan(&results).Error
	return results, err
}

func (r *productRepository) FetchSizesType(ctx context.Context, productID int, customerTypeID int) ([]domain.ProductSizeTypeResponse, error) {
	var results []domain.ProductSizeTypeResponse
	err := r.db.WithContext(ctx).Table("t_product_price").
		Select("t_product_price.id_product_price as id, t_product_size.name_product_size as size_name, t_product_price.price, t_product_price.stock, t_product_price.weight").
		Joins("LEFT JOIN t_product_size ON t_product_price.id_product_size = t_product_size.id_product_size").
		Where("t_product_price.id_product = ? AND t_product_price.id_customer_type = ?", productID, customerTypeID).
		Order("t_product_price.id_product_size ASC").
		Scan(&results).Error
	return results, err
}
