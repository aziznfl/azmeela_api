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

type ProductCodeResponse struct {
	ID            int                  `json:"id"`
	ProductTypeID int                  `json:"product_type_id"`
	Name          string               `json:"name"`
	WebStatus     int                  `json:"web_status"`
	CodeStatus    int                  `json:"code_status"`
	Description   string               `json:"description"`
	Information   string               `json:"information"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	Type          *ProductTypeResponse `json:"type,omitempty"`
	Products      []ProductSummary     `json:"products,omitempty"`
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
	ID            int     `json:"id"`
	ProductSizeID int     `json:"product_size_id"`
	SizeName      string  `json:"size_name"`
	Price         float64 `json:"price"`
	Stock         int     `json:"stock"`
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

// Mappers

func ToProductTypeResponse(t *domain.ProductType) *ProductTypeResponse {
	if t == nil {
		return nil
	}
	return &ProductTypeResponse{
		ID:        t.ID,
		Name:      t.Name,
		WebStatus: t.WebStatus,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
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
		UpdatedAt:     c.UpdatedAt,
	}

	if c.Type != nil {
		resp.Type = ToProductTypeResponse(c.Type)
	}

	if len(c.Products) > 0 {
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
					summary.Variants[j] = PriceSummary{
						ID:            v.ID,
						ProductSizeID: v.ProductSizeID,
						SizeName:      v.Size.Name,
						Price:         v.Price,
						Stock:         v.Stock,
					}
				}
			}
			resp.Products[i] = summary
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
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.ProductCode != nil {
		resp.ProductCode = ToProductCodeResponse(p.ProductCode)
	}

	if len(p.Variants) > 0 {
		resp.Variants = make([]PriceSummary, len(p.Variants))
		for i, v := range p.Variants {
			sizeName := ""
			if v.Size != nil {
				sizeName = v.Size.Name
			}
			resp.Variants[i] = PriceSummary{
				ID:            v.ID,
				ProductSizeID: v.ProductSizeID,
				SizeName:      sizeName,
				Price:         v.Price,
				Stock:         v.Stock,
			}
		}
	}

	return resp
}
