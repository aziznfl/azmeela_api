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
	err := r.db.WithContext(ctx).Order("id ASC").Find(&types).Error
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

	// Senior Optimization: Use Joins for 1:1 and Preload for 1:N with targeted selects if possible
	query = query.Joins("Type").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Order("products.color ASC")
		})

	if custTypeID, ok := filter["customer_type_id"].(int); ok && custTypeID != 0 {
		query = query.Where("EXISTS (SELECT 1 FROM products p JOIN product_prices pp ON p.id = pp.product_id WHERE p.product_code_id = product_codes.id AND pp.customer_type_id = ?)", custTypeID)
		query = query.Preload("Products.Variants", "customer_type_id = ?", custTypeID)
	} else {
		query = query.Preload("Products.Variants")
	}
	
	query = query.Preload("Products.Variants.Size")

	if name, ok := filter["name"].(string); ok && name != "" {
		query = query.Where("product_codes.name ILIKE ?", "%"+name+"%")
	}
	
	if typeID, ok := filter["product_type_id"].(int); ok && typeID != 0 {
		query = query.Where("product_codes.product_type_id = ?", typeID)
	}
	
	if productID, ok := filter["product_id"].(int); ok && productID != 0 {
		// This filter seems to refer to a specific product within a code
		query = query.Where("EXISTS (SELECT 1 FROM products p WHERE p.product_code_id = product_codes.id AND p.id = ?)", productID)
	}

	err := query.Order("product_codes.updated_at DESC").Find(&codes).Error
	return codes, err
}

func (r *productRepository) GetCodeByID(ctx context.Context, id int) (*domain.ProductCode, error) {
	var code domain.ProductCode
	err := r.db.WithContext(ctx).
		Joins("Type").
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
		query = query.Where("product_prices.product_id = ?", productID)
	}

	if customerTypeID, ok := filter["customer_type_id"].(int); ok && customerTypeID != 0 {
		query = query.Where("product_prices.customer_type_id = ?", customerTypeID)
	}

	err := query.Order("product_prices.updated_at DESC").Find(&prices).Error
	return prices, err
}

func (r *productRepository) UpdateStock(ctx context.Context, id int, stock int) error {
	return r.db.WithContext(ctx).Model(&domain.ProductPrice{}).Where("id = ?", id).Update("stock", stock).Error
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
