package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type customerUsecase struct {
	customerRepo domain.CustomerRepository
	redisRepo    domain.RedisRepository
}

func NewCustomerUsecase(c domain.CustomerRepository, redisRepo domain.RedisRepository) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo: c,
		redisRepo:    redisRepo,
	}
}

type customerFetchResponse struct {
	Customers []domain.Customer     `json:"customers"`
	Meta      domain.PaginationMeta `json:"meta"`
}

func (u *customerUsecase) Fetch(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.Customer, domain.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("customer_list:%v:%d:%d", filter, page, limit)
	var cached customerFetchResponse
	if err := u.redisRepo.Get(ctx, cacheKey, &cached); err == nil {
		return cached.Customers, cached.Meta, nil
	}

	offset := (page - 1) * limit
	customers, total, err := u.customerRepo.Fetch(ctx, filter, offset, limit)
	if err != nil {
		return nil, domain.PaginationMeta{}, err
	}

	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}

	meta := domain.PaginationMeta{
		Total:       total,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}

	_ = u.redisRepo.Set(ctx, cacheKey, customerFetchResponse{Customers: customers, Meta: meta}, 30*time.Minute)

	return customers, meta, nil
}

func (u *customerUsecase) GetByID(ctx context.Context, id int) (*domain.Customer, error) {
	return u.customerRepo.GetByID(ctx, id)
}

func (u *customerUsecase) Create(ctx context.Context, req *domain.CustomerRequest) (*domain.Customer, error) {
	// Optional: Check if username already exists
	if req.Username != nil {
		existing, _ := u.customerRepo.GetByUsername(ctx, *req.Username)
		if existing != nil {
			return nil, errors.New("username already exists")
		}
	}

	cust := &domain.Customer{
		CustomerTypeID:   req.CustomerTypeID,
		FullName:         req.FullName,
		PhoneNumber:      req.PhoneNumber,
		Deposit:          req.Deposit,
		MembershipStatus: req.MembershipStatus,
		Username:         req.Username,
		Email:            req.Email,
	}

	// Default membership status if empty
	if cust.MembershipStatus == "" {
		cust.MembershipStatus = "1"
	}

	err := u.customerRepo.Store(ctx, cust)
	if err != nil {
		return nil, err
	}

	u.clearCustomerCache(ctx)

	// Reload for preloading type
	return u.customerRepo.GetByID(ctx, cust.ID)
}

func (u *customerUsecase) Update(ctx context.Context, id int, req *domain.CustomerRequest) (*domain.Customer, error) {
	cust, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cust.CustomerTypeID = req.CustomerTypeID
	cust.FullName = req.FullName
	cust.PhoneNumber = req.PhoneNumber
	cust.Deposit = req.Deposit
	cust.MembershipStatus = req.MembershipStatus
	cust.Username = req.Username
	cust.Email = req.Email

	err = u.customerRepo.Update(ctx, cust)
	if err != nil {
		return nil, err
	}

	u.clearCustomerCache(ctx)

	return u.customerRepo.GetByID(ctx, cust.ID)
}

func (u *customerUsecase) Delete(ctx context.Context, id int) error {
	err := u.customerRepo.Delete(ctx, id)
	if err == nil {
		u.clearCustomerCache(ctx)
	}
	return err
}

func (u *customerUsecase) GetTypes(ctx context.Context) ([]domain.CustomerType, error) {
	cacheKey := "customer_types"
	var types []domain.CustomerType
	if err := u.redisRepo.Get(ctx, cacheKey, &types); err == nil {
		return types, nil
	}

	types, err := u.customerRepo.FetchTypes(ctx)
	if err != nil {
		return nil, err
	}

	_ = u.redisRepo.Set(ctx, cacheKey, types, 24*time.Hour)
	return types, nil
}

func (u *customerUsecase) CreateType(ctx context.Context, req *domain.CustomerType) error {
	err := u.customerRepo.CreateType(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "customer_types")
	}
	return err
}

func (u *customerUsecase) UpdateType(ctx context.Context, id int, req *domain.CustomerType) error {
	req.ID = id
	err := u.customerRepo.UpdateType(ctx, req)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "customer_types")
	}
	return err
}

func (u *customerUsecase) DeleteType(ctx context.Context, id int) error {
	err := u.customerRepo.DeleteType(ctx, id)
	if err == nil {
		_ = u.redisRepo.Delete(ctx, "customer_types")
	}
	return err
}

func (u *customerUsecase) GetAddresses(ctx context.Context, customerID int) ([]domain.CustomerAddress, error) {
	return u.customerRepo.FetchAddresses(ctx, customerID)
}

func (u *customerUsecase) CreateAddress(ctx context.Context, req *domain.CustomerAddressRequest) (*domain.CustomerAddress, error) {
	addr := &domain.CustomerAddress{
		CustomerID:    req.CustomerID,
		Country:       req.Country,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		SubDistrict:   req.SubDistrict,
		StreetAddress: req.StreetAddress,
		PostalCode:    req.PostalCode,
		SicepatID:     req.SicepatID,
	}

	err := u.customerRepo.StoreAddress(ctx, addr)
	if err != nil {
		return nil, err
	}

	return u.customerRepo.GetAddressByID(ctx, addr.ID)
}

func (u *customerUsecase) UpdateAddress(ctx context.Context, id int, req *domain.CustomerAddressRequest) (*domain.CustomerAddress, error) {
	addr, err := u.customerRepo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	addr.Country = req.Country
	addr.Province = req.Province
	addr.City = req.City
	addr.District = req.District
	addr.SubDistrict = req.SubDistrict
	addr.StreetAddress = req.StreetAddress
	addr.PostalCode = req.PostalCode
	addr.SicepatID = req.SicepatID

	err = u.customerRepo.UpdateAddress(ctx, addr)
	if err != nil {
		return nil, err
	}

	return u.customerRepo.GetAddressByID(ctx, addr.ID)
}

func (u *customerUsecase) DeleteAddress(ctx context.Context, id int) error {
	return u.customerRepo.DeleteAddress(ctx, id)
}

func (u *customerUsecase) clearCustomerCache(ctx context.Context) {
	// Simplified cache clearing
	// In production, use Redis SCAN or pattern match if keys are many
}
