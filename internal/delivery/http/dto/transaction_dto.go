package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type TransactionResponse struct {
	ID              int                         `json:"id"`
	TransactionCode string                      `json:"transaction_code"`
	Total           int                         `json:"total"`
	TotalPrice      int                         `json:"totalPrice"`
	CreatedAt       time.Time                   `json:"created_at"`
	PaymentCode     int                         `json:"payment_code,omitempty"`
	IsReminded      bool                        `json:"is_reminded,omitempty"`
	TransferDate    *time.Time                  `json:"transfer_date,omitempty"`
	ShippingDate    *time.Time                  `json:"shipping_date,omitempty"`
	Customer        *TransactionCustomerSummary `json:"customer,omitempty"`
	Status          *TransactionStatusSummary   `json:"status,omitempty"`
	Shipping        *TransactionShippingSummary `json:"shipping,omitempty"`
	Discount        *TransactionDiscountSummary `json:"discount,omitempty"`
	Details         []TransactionDetailResponse `json:"details,omitempty"`
}

type TransactionCustomerSummary struct {
	ID               int    `json:"id"`
	FullName         string `json:"full_name"`
	CustomerTypeID   int    `json:"customer_type_id"`
	CustomerTypeName string `json:"customer_type_name"`
	MembershipStatus string `json:"membership_status"`
	PhoneNumber      string `json:"phone_number"`
}

type TransactionStatusSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TransactionShippingSummary struct {
	ShippingCost   int     `json:"shipping_cost"`
	TrackingNumber *string `json:"tracking_number"`
	Courier        *string `json:"courier"`
	Address        string  `json:"address"`
}

type TransactionDiscountSummary struct {
	Discount     int    `json:"discount"`
	DiscountNote string `json:"discount_note"`
	DiscountType int    `json:"discount_type"`
}

type TransactionDetailResponse struct {
	ID              int    `json:"id"`
	ProductPriceID  int    `json:"product_price_id"`
	ProductCodeName string `json:"product_code_name"`
	ProductTypeName string `json:"product_type_name"`
	SKU             string `json:"sku,omitempty"`
	Color           string `json:"color"`
	Size            string `json:"size"`
	Quantity        int    `json:"quantity"`
	Price           int    `json:"price"`
	TotalPrice      int    `json:"total_price"`
}

type TransactionLogResponse struct {
	ID        int                       `json:"id"`
	Admin     *TransactionAdminSummary  `json:"admin"`
	OldStatus *TransactionStatusSummary `json:"old_status"`
	NewStatus *TransactionStatusSummary `json:"new_status"`
	CreatedAt time.Time                 `json:"created_at"`
}

type TransactionRequest struct {
	CustomerID      int                        `json:"customer_id" binding:"required"`
	AdminID         int                        `json:"admin_id" binding:"required"`
	TransferDate    *time.Time                 `json:"transfer_date"`
	ShippingDate    *time.Time                 `json:"shipping_date"`
	StatusID        int                        `json:"status_id"`
	ShippingCost    int                        `json:"shipping_cost"`
	TrackingNumber  *string                    `json:"tracking_number"`
	Courier         *string                    `json:"courier"`
	TransactionCode string                     `json:"transaction_code" binding:"required"`
	Total           int                        `json:"total" binding:"required"`
	Address         string                     `json:"address" binding:"required"`
	PaymentCode     int                        `json:"payment_code"`
	Discount        int                        `json:"discount"`
	DiscountNote    *string                    `json:"discount_note"`
	DiscountType    int                        `json:"discount_type"`
	IsReminded      bool                       `json:"is_reminded"`
	Details         []TransactionDetailRequest `json:"details"`
}

type TransactionDetailRequest struct {
	ProductPriceID int `json:"product_price_id" binding:"required"`
	Quantity       int `json:"quantity" binding:"required,gt=0"`
	Price          int `json:"price" binding:"required"`
}

type TransactionAdminSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ToTransactionResponse(t *domain.Transaction) *TransactionResponse {
	if t == nil {
		return nil
	}

	resp := &TransactionResponse{
		ID:              t.ID,
		TransactionCode: t.TransactionCode,
		Total:           t.Total,
		PaymentCode:     t.PaymentCode,
		IsReminded:      t.IsReminded,
		TransferDate:    t.TransferDate,
		ShippingDate:    t.ShippingDate,
		CreatedAt:       t.CreatedAt,
	}

	if t.Customer != nil {
		phone := ""
		if t.Customer.PhoneNumber != nil {
			phone = *t.Customer.PhoneNumber
		}
		resp.Customer = &TransactionCustomerSummary{
			ID:               t.Customer.ID,
			FullName:         t.Customer.FullName,
			CustomerTypeID:   t.Customer.CustomerTypeID,
			CustomerTypeName: t.Customer.CustomerType.Name,
			MembershipStatus: t.Customer.MembershipStatus,
			PhoneNumber:      phone,
		}
	}

	if t.Status != nil {
		resp.Status = &TransactionStatusSummary{
			ID:   t.Status.ID,
			Name: t.Status.Name,
		}
	}

	resp.Shipping = &TransactionShippingSummary{
		ShippingCost:   t.ShippingCost,
		TrackingNumber: t.TrackingNumber,
		Courier:        t.Courier,
		Address:        t.Address,
	}

	resp.Discount = &TransactionDiscountSummary{
		Discount:     t.Discount,
		DiscountNote: "",
		DiscountType: t.DiscountType,
	}
	if t.DiscountNote != nil {
		resp.Discount.DiscountNote = *t.DiscountNote
	}

	if len(t.Details) > 0 {
		resp.Details = make([]TransactionDetailResponse, len(t.Details))
		for i, d := range t.Details {
			detail := TransactionDetailResponse{
				ID:             d.ID,
				ProductPriceID: d.ProductPriceID,
				Quantity:       d.Quantity,
				Price:          d.Price,
				TotalPrice:     d.Quantity * d.Price,
			}
			if d.ProductPrice != nil {
				if d.ProductPrice.Product != nil {
					p := d.ProductPrice.Product
					detail.SKU = p.SKU
					detail.Color = p.Color
					if p.ProductCode != nil {
						detail.ProductCodeName = p.ProductCode.Name
						if p.ProductCode.Type != nil {
							detail.ProductTypeName = p.ProductCode.Type.Name
						}
					}
				}
				if d.ProductPrice.Size != nil {
					detail.Size = d.ProductPrice.Size.Name
				}
			}
			resp.Details[i] = detail
		}
	}

	discountValue := t.Discount
	if t.DiscountType == 1 { // Percentage
		discountValue = (t.Total * t.Discount) / 100
	}

	resp.TotalPrice = t.Total + t.ShippingCost + t.PaymentCode - discountValue

	return resp
}

func ToTransactionListResponse(transactions []domain.Transaction) []*TransactionResponse {
	resps := make([]*TransactionResponse, len(transactions))
	for i, t := range transactions {
		resps[i] = ToTransactionResponse(&t)
	}
	return resps
}

func ToTransactionLogResponse(log *domain.TransactionLog) *TransactionLogResponse {
	if log == nil {
		return nil
	}
	resp := &TransactionLogResponse{
		ID:        log.ID,
		CreatedAt: log.CreatedAt,
	}
	if log.Admin != nil {
		resp.Admin = &TransactionAdminSummary{
			ID:   log.Admin.ID,
			Name: log.Admin.Name,
		}
	}
	if log.OldStatus != nil {
		resp.OldStatus = &TransactionStatusSummary{
			ID:   log.OldStatus.ID,
			Name: log.OldStatus.Name,
		}
	}
	if log.NewStatus != nil {
		resp.NewStatus = &TransactionStatusSummary{
			ID:   log.NewStatus.ID,
			Name: log.NewStatus.Name,
		}
	}
	return resp
}

func ToTransactionLogListResponse(logs []domain.TransactionLog) []*TransactionLogResponse {
	resps := make([]*TransactionLogResponse, len(logs))
	for i, l := range logs {
		resps[i] = ToTransactionLogResponse(&l)
	}
	return resps
}

func ToTransactionStatusListResponse(statuses []domain.TransactionStatus) []*TransactionStatusSummary {
	resps := make([]*TransactionStatusSummary, len(statuses))
	for i, s := range statuses {
		resps[i] = &TransactionStatusSummary{
			ID:   s.ID,
			Name: s.Name,
		}
	}
	return resps
}

func (r *TransactionRequest) ToDomain() *domain.TransactionRequest {
	details := make([]domain.TransactionDetailRequest, len(r.Details))
	for i, d := range r.Details {
		details[i] = domain.TransactionDetailRequest{
			ProductPriceID: d.ProductPriceID,
			Quantity:       d.Quantity,
			Price:          d.Price,
		}
	}

	return &domain.TransactionRequest{
		CustomerID:      r.CustomerID,
		AdminID:         r.AdminID,
		TransferDate:    r.TransferDate,
		ShippingDate:    r.ShippingDate,
		StatusID:        r.StatusID,
		ShippingCost:    r.ShippingCost,
		TrackingNumber:  r.TrackingNumber,
		Courier:         r.Courier,
		TransactionCode: r.TransactionCode,
		Total:           r.Total,
		Address:         r.Address,
		PaymentCode:     r.PaymentCode,
		Discount:        r.Discount,
		DiscountNote:    r.DiscountNote,
		DiscountType:    r.DiscountType,
		IsReminded:      r.IsReminded,
		Details:         details,
	}
}
