package dto

import (
	"github.com/azmeela/sispeg-api/internal/domain"
)

type PayrollSummaryResponse struct {
	EmployeeID                int                           `json:"employee_id"`
	EmployeeName              string                        `json:"employee_name"`
	BaseSalary                int                           `json:"base_salary"`
	OvertimeInfo              []OvertimePayrollResponse     `json:"overtime_info"`
	OvertimePay               int                           `json:"overtime_pay"`
	CashAdvances              []CashAdvanceDeductionResponse `json:"cash_advances"`
	TotalDeductionCashAdvance int                           `json:"total_deduction_cash_advance"`
	SalaryComponents          []SalaryComponentLineResponse `json:"salary_components"`
	TotalAllowance            int                           `json:"total_allowance"`
	TotalDeduction            int                           `json:"total_deduction"`
	NetSalary                 int                           `json:"net_salary"`
}

type OvertimePayrollResponse struct {
	Date        string  `json:"date"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Hours       float64 `json:"hours"`
	RatePerHour int     `json:"rate_per_hour"`
	Total       int     `json:"total"`
}

type CashAdvanceDeductionResponse struct {
	ID        int    `json:"id"`
	Amount    int    `json:"amount"`
	Purpose   string `json:"purpose"`
	CreatedAt string `json:"created_at"`
}

type SalaryComponentLineResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  int    `json:"type"` // 1: allowance, 2: deduction
	Value int    `json:"value"`
}

type PayrollRequest struct {
	Month int `json:"month" form:"month" binding:"required,min=1,max=12"`
	Year  int `json:"year" form:"year" binding:"required,min=2000"`
}

func ToPayrollSummaryResponse(p *domain.PayrollSummary) *PayrollSummaryResponse {
	if p == nil {
		return nil
	}

	overtimes := make([]OvertimePayrollResponse, len(p.OvertimeInfo))
	for i, o := range p.OvertimeInfo {
		overtimes[i] = OvertimePayrollResponse(o)
	}

	advances := make([]CashAdvanceDeductionResponse, len(p.CashAdvances))
	for i, a := range p.CashAdvances {
		advances[i] = CashAdvanceDeductionResponse(a)
	}

	components := make([]SalaryComponentLineResponse, len(p.SalaryComponents))
	for i, s := range p.SalaryComponents {
		components[i] = SalaryComponentLineResponse(s)
	}

	return &PayrollSummaryResponse{
		EmployeeID:                p.EmployeeID,
		EmployeeName:              p.EmployeeName,
		BaseSalary:                p.BaseSalary,
		OvertimeInfo:              overtimes,
		OvertimePay:               p.OvertimePay,
		CashAdvances:              advances,
		TotalDeductionCashAdvance: p.TotalDeductionCashAdvance,
		SalaryComponents:          components,
		TotalAllowance:            p.TotalAllowance,
		TotalDeduction:            p.TotalDeduction,
		NetSalary:                 p.NetSalary,
	}
}

func ToPayrollListResponse(items []domain.PayrollSummary) []*PayrollSummaryResponse {
	resps := make([]*PayrollSummaryResponse, len(items))
	for i, item := range items {
		resps[i] = ToPayrollSummaryResponse(&item)
	}
	return resps
}
