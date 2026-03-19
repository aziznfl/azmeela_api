package usecase

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

func (u *productUsecase) GetInventoryList(ctx context.Context, filter map[string]interface{}) ([]domain.ProductCode, error) {
	return u.productRepo.FetchCodes(ctx, filter)
}

func (u *productUsecase) GetProductTypes(ctx context.Context) ([]domain.ProductType, error) {
	return u.productRepo.FetchTypes(ctx)
}

func (u *productUsecase) CreateProductType(ctx context.Context, req *domain.ProductType) error {
	return u.productRepo.CreateType(ctx, req)
}

func (u *productUsecase) UpdateProductType(ctx context.Context, req *domain.ProductType) error {
	return u.productRepo.UpdateType(ctx, req)
}

func (u *productUsecase) DeleteProductType(ctx context.Context, id int) error {
	return u.productRepo.DeleteType(ctx, id)
}

func (u *productUsecase) CreateProductCode(ctx context.Context, req *domain.ProductCode) error {
	return u.productRepo.CreateCode(ctx, req)
}

func (u *productUsecase) UpdateProductCode(ctx context.Context, req *domain.ProductCode) error {
	return u.productRepo.UpdateCode(ctx, req)
}

func (u *productUsecase) DeleteProductCode(ctx context.Context, id int) error {
	return u.productRepo.DeleteCode(ctx, id)
}

func (u *productUsecase) GetProductSizes(ctx context.Context) ([]domain.ProductSize, error) {
	return u.productRepo.FetchSizes(ctx)
}

func (u *productUsecase) UpdateStock(ctx context.Context, id int, stock int) error {
	return u.productRepo.UpdateStock(ctx, id, stock)
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *domain.Product) error {
	return u.productRepo.CreateProduct(ctx, req)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, req *domain.Product) error {
	return u.productRepo.UpdateProduct(ctx, req)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id int) error {
	return u.productRepo.DeleteProduct(ctx, id)
}
