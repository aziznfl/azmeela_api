package repository

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) domain.CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) Fetch(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Customer{}).Where("is_active = ?", true)

	order := "DESC"
	if ord, ok := filter["order"]; ok {
		order = ord.(string)
		delete(filter, "order")
	}

	sortOrder := "customers.created_at " + order
	if sortBy, ok := filter["sort_by"]; ok {
		if sortBy == "id" {
			sortOrder = "customers.id " + order
		}
		delete(filter, "sort_by")
	}

	for k, v := range filter {
		if k == "search" {
			searchVal := "%" + v.(string) + "%"
			query = query.Where("LOWER(full_name) LIKE LOWER(?)", searchVal)
		} else {
			query = query.Where(k+" = ?", v)
		}
	}

	// Get Total Count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply Pagination & Joins (Senior Optimization: one query for 1:1 relations)
	err = query.Joins("CustomerType").Offset(offset).Limit(limit).Order(sortOrder).Find(&customers).Error

	return customers, total, err
}

func (r *customerRepository) GetByID(ctx context.Context, id int) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).Joins("CustomerType").Where("customers.is_active = ?", true).First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByUsername(ctx context.Context, username string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.WithContext(ctx).Joins("CustomerType").Where("customers.username = ? AND customers.is_active = ?", username, true).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Store(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&domain.Customer{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *customerRepository) FetchTypes(ctx context.Context) ([]domain.CustomerType, error) {
	var types []domain.CustomerType
	err := r.db.WithContext(ctx).Order("id ASC").Find(&types).Error
	return types, err
}

func (r *customerRepository) GetTypeByID(ctx context.Context, id int) (*domain.CustomerType, error) {
	var cType domain.CustomerType
	err := r.db.WithContext(ctx).First(&cType, id).Error
	return &cType, err
}

func (r *customerRepository) CreateType(ctx context.Context, cType *domain.CustomerType) error {
	return r.db.WithContext(ctx).Create(cType).Error
}

func (r *customerRepository) UpdateType(ctx context.Context, cType *domain.CustomerType) error {
	return r.db.WithContext(ctx).Save(cType).Error
}

func (r *customerRepository) DeleteType(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.CustomerType{}, id).Error
}

func (r *customerRepository) FetchAddresses(ctx context.Context, customerID int) ([]domain.CustomerAddress, error) {
	var addresses []domain.CustomerAddress
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Order("id ASC").Find(&addresses).Error
	return addresses, err
}

func (r *customerRepository) GetAddressByID(ctx context.Context, id int) (*domain.CustomerAddress, error) {
	var address domain.CustomerAddress
	err := r.db.WithContext(ctx).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *customerRepository) StoreAddress(ctx context.Context, address *domain.CustomerAddress) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *customerRepository) UpdateAddress(ctx context.Context, address *domain.CustomerAddress) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *customerRepository) DeleteAddress(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.CustomerAddress{}, id).Error
}
