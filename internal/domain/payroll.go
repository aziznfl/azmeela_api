package domain

import "context"

// PayrollSummary represents one employee's salary calculation
type PayrollSummary struct {
	EmployeeID                int
	EmployeeName              string
	BaseSalary                int
	OvertimeInfo              []OvertimePayroll
	OvertimePay               int
	CashAdvances              []CashAdvanceDeduction
	TotalDeductionCashAdvance int
	SalaryComponents          []SalaryComponentLine
	TotalAllowance            int
	TotalDeduction            int
	NetSalary                 int
}

// OvertimePayroll represents overtime data used in payroll
type OvertimePayroll struct {
	Date        string
	StartTime   string
	EndTime     string
	Hours       float64
	RatePerHour int
	Total       int
}

// CashAdvanceDeduction represents a cash advance deducted from salary
type CashAdvanceDeduction struct {
	ID        int
	Amount    int
	Purpose   string
	CreatedAt string
}

// SalaryComponentLine represents an allowance or deduction from salary_variables
type SalaryComponentLine struct {
	ID    int
	Name  string
	Type  int // 1: allowance, 2: deduction
	Value int
}

// PayrollRequest contains the filter for generating payroll
type PayrollRequest struct {
	Month int
	Year  int
}

// PayrollUsecase represents the payroll usecase contract
type PayrollUsecase interface {
	Generate(ctx context.Context, month int, year int) ([]PayrollSummary, error)
	GenerateByEmployee(ctx context.Context, employeeID int, month int, year int) (*PayrollSummary, error)
}
