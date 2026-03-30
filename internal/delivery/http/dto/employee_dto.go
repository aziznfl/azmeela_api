package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type EmployeeResponse struct {
	ID            int        `json:"id"`
	TypeID        int        `json:"type_id"`
	Username      string     `json:"username"`
	Name          string     `json:"name"`
	Status        bool       `json:"status"`
	Bio           string     `json:"bio"`
	BaseSalary    int        `json:"base_salary"`
	ContractStart *time.Time `json:"contract_start"`
	ContractEnd   *time.Time `json:"contract_end"`
	CV            string     `json:"cv"`
	Phone         string     `json:"phone"`
	AdminTypeName string     `json:"admin_type_name,omitempty"`
}

type EmployeeRequest struct {
	TypeID        int        `json:"type_id" binding:"required"`
	Username      string     `json:"username" binding:"required"`
	Password      string     `json:"password,omitempty"`
	Name          string     `json:"name" binding:"required"`
	Active        bool       `json:"status"`
	Bio           string     `json:"bio"`
	BaseSalary    int        `json:"base_salary"`
	ContractStart *time.Time `json:"contract_start"`
	ContractEnd   *time.Time `json:"contract_end"`
	CV            string     `json:"cv"`
	Phone         string     `json:"phone"`
}

type AdminTypeResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
}

func ToEmployeeResponse(e *domain.Employee) *EmployeeResponse {
	if e == nil {
		return nil
	}

	resp := &EmployeeResponse{
		ID:            e.ID,
		TypeID:        e.TypeID,
		Username:      e.Username,
		Name:          e.Name,
		Status:        e.Active == 1,
		Bio:           e.Bio,
		BaseSalary:    e.BaseSalary,
		ContractStart: e.ContractStart,
		ContractEnd:   e.ContractEnd,
		CV:            e.CV,
		Phone:         e.Phone,
	}

	if e.AdminType != nil {
		resp.AdminTypeName = e.AdminType.Name
	}

	return resp
}

func ToEmployeeListResponse(employees []domain.Employee) []*EmployeeResponse {
	resps := make([]*EmployeeResponse, len(employees))
	for i, e := range employees {
		resps[i] = ToEmployeeResponse(&e)
	}
	return resps
}

func ToAdminTypeResponse(at *domain.AdminType) *AdminTypeResponse {
	if at == nil {
		return nil
	}
	return &AdminTypeResponse{
		ID:        at.ID,
		Name:      at.Name,
	}
}

func ToAdminTypeListResponse(types []domain.AdminType) []*AdminTypeResponse {
	resps := make([]*AdminTypeResponse, len(types))
	for i, t := range types {
		resps[i] = ToAdminTypeResponse(&t)
	}
	return resps
}
func (r *EmployeeRequest) ToDomain() *domain.Employee {
	active := 0
	if r.Active {
		active = 1
	}
	return &domain.Employee{
		TypeID:        r.TypeID,
		Username:      r.Username,
		Password:      r.Password,
		Name:          r.Name,
		Active:        active,
		Bio:           r.Bio,
		BaseSalary:    r.BaseSalary,
		ContractStart: r.ContractStart,
		ContractEnd:   r.ContractEnd,
		CV:            r.CV,
		Phone:         r.Phone,
	}
}
