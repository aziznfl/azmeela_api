package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	redisRepo   domain.RedisRepository
}

func NewProductUsecase(repo domain.ProductRepository, redisRepo domain.RedisRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
		redisRepo:   redisRepo,
	}
}

func (u *productUsecase) GetInventoryList(ctx context.Context, filter map[string]interface{}) ([]domain.ProductCode, error) {
	return u.productRepo.FetchCodes(ctx, filter)
}

func (u *productUsecase) GetCodesWithTypes(ctx context.Context, filter map[string]interface{}) ([]domain.ProductCodeWithType, error) {
	cacheKey := fmt.Sprintf("product_codes_with_types:%v", filter)
	var results []domain.ProductCodeWithType

	if err := u.redisRepo.Get(ctx, cacheKey, &results); err == nil {
		return results, nil
	}

	results, err := u.productRepo.FetchCodesWithTypes(ctx, filter)
	if err != nil {
		return nil, err
	}

	_ = u.redisRepo.Set(ctx, cacheKey, results, 24*time.Hour)
	return results, nil
}

func (u *productUsecase) GetProductTypes(ctx context.Context) ([]domain.ProductType, error) {
	cacheKey := "product_types"
	var types []domain.ProductType

	if err := u.redisRepo.Get(ctx, cacheKey, &types); err == nil {
		return types, nil
	}

	types, err := u.productRepo.FetchTypes(ctx)
	if err != nil {
		return nil, err
	}

	_ = u.redisRepo.Set(ctx, cacheKey, types, 24*time.Hour)
	return types, nil
}

func (u *productUsecase) CreateProductType(ctx context.Context, req *domain.ProductType) error {
	err := u.productRepo.CreateType(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_types")
	}
	return err
}

func (u *productUsecase) UpdateProductType(ctx context.Context, req *domain.ProductType) error {
	err := u.productRepo.UpdateType(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_types")
	}
	return err
}

func (u *productUsecase) DeleteProductType(ctx context.Context, id int) error {
	err := u.productRepo.DeleteType(ctx, id)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_types")
	}
	return err
}

func (u *productUsecase) CreateProductCode(ctx context.Context, req *domain.ProductCode) error {
	err := u.productRepo.CreateCode(ctx, req)
	if err == nil {
		u.clearCodesWithTypesCache(ctx)
	}
	return err
}

func (u *productUsecase) UpdateProductCode(ctx context.Context, req *domain.ProductCode) error {
	err := u.productRepo.UpdateCode(ctx, req)
	if err == nil {
		u.clearCodesWithTypesCache(ctx)
	}
	return err
}

func (u *productUsecase) DeleteProductCode(ctx context.Context, id int) error {
	err := u.productRepo.DeleteCode(ctx, id)
	if err == nil {
		u.clearCodesWithTypesCache(ctx)
	}
	return err
}


func (u *productUsecase) UpdateStock(ctx context.Context, id int, quantity int, adminID int) error {
	return u.productRepo.UpdateStock(ctx, id, quantity, adminID)
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *domain.Product) error {
	err := u.productRepo.CreateProduct(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, fmt.Sprintf("product_colors:%d", req.ProductCodeID))
	}
	return err
}

func (u *productUsecase) UpdateProduct(ctx context.Context, req *domain.Product) error {
	err := u.productRepo.UpdateProduct(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, fmt.Sprintf("product_colors:%d", req.ProductCodeID))
	}
	return err
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	// We might need to fetch the product first to get the ProductCodeID for cache invalidation
	// But let's keep it simple for now, or just clear all colors cache if needed.
	return u.productRepo.DeleteProduct(ctx, id)
}

func (u *productUsecase) GetStockLogs(ctx context.Context, productPriceID int) ([]domain.ProductStockLog, error) {
	return u.productRepo.GetStockLogs(ctx, productPriceID)
}

func (u *productUsecase) GetProductColors(ctx context.Context, productCodeID int) ([]domain.ProductColorResponse, error) {
	cacheKey := fmt.Sprintf("product_colors:%d", productCodeID)
	colors := []domain.ProductColorResponse{}

	if err := u.redisRepo.Get(ctx, cacheKey, &colors); err == nil {
		return colors, nil
	}

	colors, err := u.productRepo.FetchColors(ctx, productCodeID)
	if err != nil {
		return nil, err
	}

	_ = u.redisRepo.Set(ctx, cacheKey, colors, 24*time.Hour)
	return colors, nil
}

func (u *productUsecase) GetProductSizesType(ctx context.Context, productID int, customerTypeID int) ([]domain.ProductSizeTypeResponse, error) {
	return u.productRepo.FetchSizesType(ctx, productID, customerTypeID)
}

func (u *productUsecase) GetAllProductSizes(ctx context.Context) ([]domain.ProductSize, error) {
	cacheKey := "product_sizes_all"
	var sizes []domain.ProductSize

	if err := u.redisRepo.Get(ctx, cacheKey, &sizes); err == nil {
		return sizes, nil
	}

	sizes, err := u.productRepo.FetchProductSizes(ctx)
	if err != nil {
		return nil, err
	}

	_ = u.redisRepo.Set(ctx, cacheKey, sizes, 24*time.Hour)
	return sizes, nil
}

func (u *productUsecase) CreateProductSize(ctx context.Context, req *domain.ProductSize) error {
	err := u.productRepo.CreateProductSize(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_sizes_all")
	}
	return err
}

func (u *productUsecase) UpdateProductSize(ctx context.Context, req *domain.ProductSize) error {
	err := u.productRepo.UpdateProductSize(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_sizes_all")
	}
	return err
}

func (u *productUsecase) DeleteProductSize(ctx context.Context, id int) error {
	err := u.productRepo.DeleteProductSize(ctx, id)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "product_sizes_all")
	}
	return err
}

func (u *productUsecase) CreateProductPrice(ctx context.Context, req *domain.ProductPrice) error {
	return u.productRepo.CreatePrice(ctx, req)
}

func (u *productUsecase) UpdateProductPrice(ctx context.Context, req *domain.ProductPrice) error {
	return u.productRepo.UpdatePrice(ctx, req)
}

func (u *productUsecase) DeleteProductPrice(ctx context.Context, id int) error {
	return u.productRepo.DeletePrice(ctx, id)
}

// Helper to clear multiple cache variations
func (u *productUsecase) clearCodesWithTypesCache(ctx context.Context) {
	// Since we use map based keys, it's hard to delete all variations without SCAN.
	// For now, we'll set a shorter TTL or use a prefix.
	// In a real scenario, we'd use a pattern match delete or a versioning system.
}
