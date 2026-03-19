package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type CustomerResponse struct {
	ID               int                       `json:"id"`
	CustomerTypeID   int                       `json:"customer_type_id"`
	CustomerType     *CustomerTypeResponse     `json:"customer_type,omitempty"`
	CustomerTypeName string                    `json:"customer_type_name,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	FullName         string                    `json:"full_name"`
	PhoneNumber      string                    `json:"phone_number"`
	Deposit          float64                   `json:"deposit"`
	MembershipStatus string                    `json:"membership_status"`
	Email            string                    `json:"email"`
	Addresses        []CustomerAddressResponse `json:"addresses,omitempty"`
}

type CustomerTypeResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Initial string `json:"initial"`
}

type CustomerAddressResponse struct {
	ID            int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	Country       string    `json:"country"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	SubDistrict   string    `json:"sub_district,omitempty"`
	StreetAddress string    `json:"street_address"`
	PostalCode    string    `json:"postal_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToCustomerResponse(c *domain.Customer) *CustomerResponse {
	if c == nil {
		return nil
	}

	phone := ""
	if c.PhoneNumber != nil {
		phone = *c.PhoneNumber
	}

	email := ""
	if c.Email != nil {
		email = *c.Email
	}

	resp := &CustomerResponse{
		ID:               c.ID,
		CustomerTypeID:   c.CustomerTypeID,
		FullName:         c.FullName,
		PhoneNumber:      phone,
		Deposit:          c.Deposit,
		MembershipStatus: c.MembershipStatus,
		Email:            email,
		CreatedAt:        c.CreatedAt,
	}

	// Map Customer Type
	resp.CustomerType = &CustomerTypeResponse{
		ID:      c.CustomerType.ID,
		Name:    c.CustomerType.Name,
		Initial: c.CustomerType.Initial,
	}
	resp.CustomerTypeName = c.CustomerType.Name

	// Map Addresses if loaded
	if len(c.Addresses) > 0 {
		resp.Addresses = make([]CustomerAddressResponse, len(c.Addresses))
		for i, a := range c.Addresses {
			resp.Addresses[i] = *ToCustomerAddressResponse(&a)
		}
	}

	return resp
}

func ToCustomerAddressResponse(a *domain.CustomerAddress) *CustomerAddressResponse {
	if a == nil {
		return nil
	}

	subDistrict := ""
	if a.SubDistrict != nil {
		subDistrict = *a.SubDistrict
	}

	postalCode := ""
	if a.PostalCode != nil {
		postalCode = *a.PostalCode
	}

	return &CustomerAddressResponse{
		ID:            a.ID,
		CustomerID:    a.CustomerID,
		Country:       a.Country,
		Province:      a.Province,
		City:          a.City,
		District:      a.District,
		SubDistrict:   subDistrict,
		StreetAddress: a.StreetAddress,
		PostalCode:    postalCode,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func ToCustomerListResponse(customers []domain.Customer) []*CustomerResponse {
	resps := make([]*CustomerResponse, len(customers))
	for i, c := range customers {
		resps[i] = ToCustomerResponse(&c)
	}
	return resps
}

func ToCustomerAddressListResponse(addresses []domain.CustomerAddress) []*CustomerAddressResponse {
	resps := make([]*CustomerAddressResponse, len(addresses))
	for i, a := range addresses {
		resps[i] = ToCustomerAddressResponse(&a)
	}
	return resps
}

func ToCustomerTypeResponse(ct *domain.CustomerType) *CustomerTypeResponse {
	if ct == nil {
		return nil
	}
	return &CustomerTypeResponse{
		ID:      ct.ID,
		Name:    ct.Name,
		Initial: ct.Initial,
	}
}

func ToCustomerTypeListResponse(types []domain.CustomerType) []*CustomerTypeResponse {
	resps := make([]*CustomerTypeResponse, len(types))
	for i, t := range types {
		resps[i] = ToCustomerTypeResponse(&t)
	}
	return resps
}
