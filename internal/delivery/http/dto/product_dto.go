package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type ProductTypeResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	WebStatus int       `json:"web_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductSizeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductSizeWithCustomerTypeResponse struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CustomerTypeID   int    `json:"customer_type_id"`
	CustomerTypeName string `json:"customer_type_name"`
}

type ProductCodeResponse struct {
	ID            int       `json:"id"`
	ProductTypeID int       `json:"product_type_id"`
	Name          string    `json:"name"`
	WebStatus     int       `json:"web_status"`
	CodeStatus    int       `json:"code_status"`
	Description   string    `json:"description"`
	Information   string               `json:"information"`
	CreatedAt     time.Time            `json:"created_at"`
	Products      []ProductSummary     `json:"products,omitempty"`
	ProductSKUs   []ProductSKUGroup    `json:"product_skus,omitempty"`
}

type ProductCodeWithTypeResponse struct {
	ID              int    `json:"id"`
	ProductTypeID   int    `json:"product_type_id"`
	Name            string `json:"name"`
	ProductTypeName string `json:"product_type_name"`
}

type ProductSKUGroup struct {
	CustomerTypeID   int              `json:"customer_type_id"`
	CustomerTypeName string           `json:"customer_type_name"`
	Products         []ProductSummary `json:"products"`
}

type ProductSummary struct {
	ID        int            `json:"id"`
	SKU       string         `json:"sku"`
	Color     string         `json:"color"`
	Status    int            `json:"status"`
	WebStatus int            `json:"web_status"`
	Variants  []PriceSummary `json:"variants,omitempty"`
}

type PriceSummary struct {
	ID              int                  `json:"id"`
	ProductID       int                  `json:"product_id"`
	ProductSizeID   int                  `json:"product_size_id"`
	CustomerTypeID  int                  `json:"customer_type_id"`
	Price           float64              `json:"price"`
	Stock           int                  `json:"stock"`
	Weight          int                  `json:"weight"`
	ProductDiscount int                  `json:"product_discount"`
	Size            *ProductSizeResponse `json:"size,omitempty"`
}

type ProductResponse struct {
	ID            int                  `json:"id"`
	ProductCodeID int                  `json:"product_code_id"`
	Status        int                  `json:"status"`
	SKU           string               `json:"sku"`
	Color         string               `json:"color"`
	Tags          string               `json:"tags"`
	WebStatus     int                  `json:"web_status"`
	SEOLink       string               `json:"seo_link"`
	Views         int                  `json:"views"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	ProductCode   *ProductCodeResponse `json:"product_code,omitempty"`
	Variants      []PriceSummary       `json:"variants,omitempty"`
}

// Requests

type CreateProductTypeRequest struct {
	Name      string `json:"name" binding:"required"`
	WebStatus int    `json:"web_status"`
}

func (r *CreateProductTypeRequest) ToEntity() *domain.ProductType {
	return &domain.ProductType{
		Name:      r.Name,
		WebStatus: r.WebStatus,
	}
}

type CreateProductCodeRequest struct {
	ProductTypeID int    `json:"product_type_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	WebStatus     int    `json:"web_status"`
	CodeStatus    int    `json:"code_status"`
	Description   string `json:"description"`
	Information   string `json:"information"`
}

func (r *CreateProductCodeRequest) ToEntity() *domain.ProductCode {
	return &domain.ProductCode{
		ProductTypeID: r.ProductTypeID,
		Name:          r.Name,
		WebStatus:     r.WebStatus,
		CodeStatus:    r.CodeStatus,
		Description:   r.Description,
		Information:   r.Information,
	}
}

type CreateProductRequest struct {
	ProductCodeID int    `json:"product_code_id" binding:"required"`
	SKU           string `json:"sku" binding:"required"`
	Color         string `json:"color" binding:"required"`
	Status        int    `json:"status"`
	WebStatus     int    `json:"web_status"`
	Tags          string `json:"tags"`
	SEOLink       string `json:"seo_link"`
}

func (r *CreateProductRequest) ToEntity() *domain.Product {
	return &domain.Product{
		ProductCodeID: r.ProductCodeID,
		SKU:           r.SKU,
		Color:         r.Color,
		Status:        r.Status,
		WebStatus:     r.WebStatus,
		Tags:          r.Tags,
		SEOLink:       r.SEOLink,
	}
}

type CreateProductSizeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (r *CreateProductSizeRequest) ToEntity() *domain.ProductSize {
	return &domain.ProductSize{
		Name: r.Name,
	}
}

type CreateProductPriceRequest struct {
	ProductID      int     `json:"product_id"`
	CustomerTypeID int     `json:"customer_type_id"`
	ProductSizeID  int     `json:"product_size_id"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	Weight         int     `json:"weight"`
}


func (r *CreateProductPriceRequest) ToEntity() *domain.ProductPrice {
	return &domain.ProductPrice{
		ProductID:      r.ProductID,
		CustomerTypeID: r.CustomerTypeID,
		ProductSizeID:  r.ProductSizeID,
		Price:          r.Price,
		Stock:          r.Stock,
		Weight:         r.Weight,
	}
}


// Mappers

func ToProductTypeResponse(t *domain.ProductType) *ProductTypeResponse {
	if t == nil {
		return nil
	}
	return &ProductTypeResponse{
		ID:        t.ID,
		Name:      t.Name,
		WebStatus: t.WebStatus,
	}
}

func ToProductTypeListResponse(types []domain.ProductType) []*ProductTypeResponse {
	resps := make([]*ProductTypeResponse, len(types))
	for i, t := range types {
		resps[i] = ToProductTypeResponse(&t)
	}
	return resps
}

func ToProductSizeResponse(s *domain.ProductSize) *ProductSizeResponse {
	if s == nil {
		return nil
	}
	return &ProductSizeResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

func ToProductSizeListResponse(sizes []domain.ProductSize) []*ProductSizeResponse {
	resps := make([]*ProductSizeResponse, len(sizes))
	for i, s := range sizes {
		resps[i] = ToProductSizeResponse(&s)
	}
	return resps
}

func ToProductSizeWithCustomerTypeResponse(s *domain.ProductSizeWithCustomerType) *ProductSizeWithCustomerTypeResponse {
	if s == nil {
		return nil
	}
	return &ProductSizeWithCustomerTypeResponse{
		ID:               s.ID,
		Name:             s.Name,
		CustomerTypeID:   s.CustomerTypeID,
		CustomerTypeName: s.CustomerTypeName,
	}
}

func ToProductSizeWithCustomerTypeListResponse(sizes []domain.ProductSizeWithCustomerType) []*ProductSizeWithCustomerTypeResponse {
	resps := make([]*ProductSizeWithCustomerTypeResponse, len(sizes))
	for i, s := range sizes {
		resps[i] = ToProductSizeWithCustomerTypeResponse(&s)
	}
	return resps
}

func ToProductCodeResponse(c *domain.ProductCode) *ProductCodeResponse {
	if c == nil {
		return nil
	}

	resp := &ProductCodeResponse{
		ID:            c.ID,
		ProductTypeID: c.ProductTypeID,
		Name:          c.Name,
		WebStatus:     c.WebStatus,
		CodeStatus:    c.CodeStatus,
		Description:   c.Description,
		Information:   c.Information,
		CreatedAt:     c.CreatedAt,
	}

	if len(c.Products) > 0 {
		skuGroups := make(map[int]*ProductSKUGroup)
		resp.Products = make([]ProductSummary, len(c.Products))

		for i, p := range c.Products {
			summary := ProductSummary{
				ID:        p.ID,
				SKU:       p.SKU,
				Color:     p.Color,
				Status:    p.Status,
				WebStatus: p.WebStatus,
			}

			if len(p.Variants) > 0 {
				summary.Variants = make([]PriceSummary, len(p.Variants))
				for j, v := range p.Variants {
					priceResp := PriceSummary{
						ID:              v.ID,
						ProductID:       v.ProductID,
						ProductSizeID:   v.ProductSizeID,
						CustomerTypeID:  v.CustomerTypeID,
						Price:           v.Price,
						Stock:           v.Stock,
						Weight:          v.Weight,
						ProductDiscount: v.ProductDiscount,
					}
					if v.Size != nil {
						priceResp.Size = ToProductSizeResponse(v.Size)
					}
					summary.Variants[j] = priceResp

					// Grouping by Customer Type
					cTypeID := v.CustomerTypeID
					if _, ok := skuGroups[cTypeID]; !ok {
						name := ""
						if v.CustomerType != nil {
							name = v.CustomerType.Name
						}
						skuGroups[cTypeID] = &ProductSKUGroup{
							CustomerTypeID:   cTypeID,
							CustomerTypeName: name,
							Products:         []ProductSummary{},
						}
					}

					group := skuGroups[cTypeID]
					// Find if the product is already in this group
					productIndex := -1
					for idx, existingProd := range group.Products {
						if existingProd.ID == p.ID {
							productIndex = idx
							break
						}
					}

					if productIndex == -1 {
						// Add product to group with this specific variant
						newProductSummary := summary
						newProductSummary.Variants = []PriceSummary{priceResp}
						group.Products = append(group.Products, newProductSummary)
					} else {
						// Append variant to existing product in group
						group.Products[productIndex].Variants = append(group.Products[productIndex].Variants, priceResp)
					}
				}
			}
			resp.Products[i] = summary
		}

		// Convert map to slice
		resp.ProductSKUs = make([]ProductSKUGroup, 0, len(skuGroups))
		for _, g := range skuGroups {
			resp.ProductSKUs = append(resp.ProductSKUs, *g)
		}
	}

	return resp
}

func ToProductCodeListResponse(codes []domain.ProductCode) []*ProductCodeResponse {
	resps := make([]*ProductCodeResponse, len(codes))
	for i, c := range codes {
		resps[i] = ToProductCodeResponse(&c)
	}
	return resps
}

func ToProductCodeWithTypeResponse(c domain.ProductCodeWithType) *ProductCodeWithTypeResponse {
	return &ProductCodeWithTypeResponse{
		ID:              c.ID,
		ProductTypeID:   c.ProductTypeID,
		Name:            c.Name,
		ProductTypeName: c.ProductTypeName,
	}
}

func ToProductCodeWithTypeListResponse(codes []domain.ProductCodeWithType) []*ProductCodeWithTypeResponse {
	resps := make([]*ProductCodeWithTypeResponse, len(codes))
	for i, c := range codes {
		resps[i] = ToProductCodeWithTypeResponse(c)
	}
	return resps
}

func ToProductResponse(p *domain.Product) *ProductResponse {
	if p == nil {
		return nil
	}

	resp := &ProductResponse{
		ID:            p.ID,
		ProductCodeID: p.ProductCodeID,
		Status:        p.Status,
		SKU:           p.SKU,
		Color:         p.Color,
		Tags:          p.Tags,
		WebStatus:     p.WebStatus,
		SEOLink:       p.SEOLink,
		Views:         p.Views,
	}

	if p.ProductCode != nil {
		resp.ProductCode = ToProductCodeResponse(p.ProductCode)
	}

	if len(p.Variants) > 0 {
		resp.Variants = make([]PriceSummary, len(p.Variants))
		for i, v := range p.Variants {
			resp.Variants[i] = PriceSummary{
				ID:              v.ID,
				ProductID:       v.ProductID,
				ProductSizeID:   v.ProductSizeID,
				CustomerTypeID:  v.CustomerTypeID,
				Price:           v.Price,
				Stock:           v.Stock,
				Weight:          v.Weight,
				ProductDiscount: v.ProductDiscount,
				Size:            ToProductSizeResponse(v.Size),
			}
		}
	}

	return resp
}

type ProductStockLogResponse struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	AdminName string    `json:"admin_name"`
	SizeName  string    `json:"size_name"`
}

func ToProductStockLogResponse(l domain.ProductStockLog) *ProductStockLogResponse {
	adminName := "Sistem"
	if l.Admin != nil {
		adminName = l.Admin.Name
	}

	sizeName := ""
	if l.ProductPrice != nil && l.ProductPrice.Size != nil {
		sizeName = l.ProductPrice.Size.Name
	}

	return &ProductStockLogResponse{
		ID:        l.ID,
		Quantity:  l.Quantity,
		CreatedAt: l.CreatedAt,
		AdminName: adminName,
		SizeName:  sizeName,
	}
}

func ToProductStockLogListResponse(logs []domain.ProductStockLog) []*ProductStockLogResponse {
	resps := make([]*ProductStockLogResponse, len(logs))
	for i, l := range logs {
		resps[i] = ToProductStockLogResponse(l)
	}
	return resps
}

type ProductColorResponse struct {
	ID            int    `json:"id"`
	ProductCodeID int    `json:"product_code_id"`
	Color         string `json:"color"`
}

func ToProductColorListResponse(colors []domain.ProductColorResponse) []ProductColorResponse {
	resps := make([]ProductColorResponse, len(colors))
	for i, c := range colors {
		resps[i] = ProductColorResponse{
			ID:            c.ID,
			ProductCodeID: c.ProductCodeID,
			Color:         c.Color,
		}
	}
	return resps
}

type ProductSizeTypeResponse struct {
	ID               int     `json:"id"`
	ProductSizeID    int     `json:"product_size_id"`
	SizeName         string  `json:"size_name"`
	Price            float64 `json:"price"`
	Stock            int     `json:"stock"`
	Weight           int     `json:"weight"`
	CustomerTypeID   int     `json:"customer_type_id"`
	CustomerTypeName string  `json:"customer_type_name"`
}

func ToProductSizeTypeListResponse(sizes []domain.ProductSizeTypeResponse) []ProductSizeTypeResponse {
	resps := make([]ProductSizeTypeResponse, len(sizes))
	for i, s := range sizes {
		resps[i] = ProductSizeTypeResponse{
			ID:               s.ID,
			ProductSizeID:    s.ProductSizeID,
			SizeName:         s.SizeName,
			Price:            s.Price,
			Stock:            s.Stock,
			Weight:           s.Weight,
			CustomerTypeID:   s.CustomerTypeID,
			CustomerTypeName: s.CustomerTypeName,
		}
	}
	return resps
}

