package usecase

import (
	"context"
	"math"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type payrollUsecase struct {
	employeeRepo    domain.EmployeeRepository
	overtimeRepo    domain.OvertimeRepository
	cashAdvanceRepo domain.CashAdvanceRepository
	salaryVarRepo   domain.SalaryVariableRepository
}

func NewPayrollUsecase(
	empRepo domain.EmployeeRepository,
	otRepo domain.OvertimeRepository,
	caRepo domain.CashAdvanceRepository,
	svRepo domain.SalaryVariableRepository,
) domain.PayrollUsecase {
	return &payrollUsecase{
		employeeRepo:    empRepo,
		overtimeRepo:    otRepo,
		cashAdvanceRepo: caRepo,
		salaryVarRepo:   svRepo,
	}
}

func (u *payrollUsecase) Generate(ctx context.Context, month int, year int) ([]domain.PayrollSummary, error) {
	employees, err := u.employeeRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.PayrollSummary
	for _, emp := range employees {
		summary, err := u.buildSummary(ctx, emp, month, year)
		if err != nil {
			return nil, err
		}
		results = append(results, *summary)
	}

	return results, nil
}

func (u *payrollUsecase) GenerateByEmployee(ctx context.Context, employeeID int, month int, year int) (*domain.PayrollSummary, error) {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	return u.buildSummary(ctx, *emp, month, year)
}

func (u *payrollUsecase) buildSummary(ctx context.Context, emp domain.Employee, month int, year int) (*domain.PayrollSummary, error) {
	// 1. Get ALL salary variables once to find the overtime rate
	allVars, err := u.salaryVarRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	// Find overtime rate from all variables
	overtimeRate := 0
	for _, sv := range allVars {
		if sv.Name == "Lembur" || sv.Name == "lembur" || sv.Name == "Overtime" {
			overtimeRate = sv.Value
			break
		}
	}

	// 2. Get approved overtimes for this employee in the given month/year
	overtimeFilter := map[string]interface{}{
		"admin_id": emp.ID,
		"status":   1, // approved
	}
	overtimes, err := u.overtimeRepo.Fetch(ctx, overtimeFilter)
	if err != nil {
		return nil, err
	}

	// Filter by month/year and calculate overtime pay
	var overtimeInfos []domain.OvertimePayroll
	totalOvertimePay := 0
	for _, ot := range overtimes {
		if ot.Date.Month() == time.Month(month) && ot.Date.Year() == year {
			hours := calculateHours(ot.TimeIn, ot.TimeOut)
			total := int(math.Ceil(hours)) * overtimeRate
			overtimeInfos = append(overtimeInfos, domain.OvertimePayroll{
				Date:        ot.Date.Format("2006-01-02"),
				StartTime:   ot.TimeIn,
				EndTime:     ot.TimeOut,
				Hours:       hours,
				RatePerHour: overtimeRate,
				Total:       total,
			})
			totalOvertimePay += total
		}
	}

	// 3. Get approved cash advances for this employee in the given month/year
	caFilter := map[string]interface{}{
		"admin_id": emp.ID,
		"status":   1, // approved
	}
	cashAdvances, err := u.cashAdvanceRepo.Fetch(ctx, caFilter)
	if err != nil {
		return nil, err
	}

	var caDeductions []domain.CashAdvanceDeduction
	totalCaDeduction := 0
	for _, ca := range cashAdvances {
		if ca.CreatedAt.Month() == time.Month(month) && ca.CreatedAt.Year() == year {
			caDeductions = append(caDeductions, domain.CashAdvanceDeduction{
				ID:        ca.ID,
				Amount:    ca.Amount,
				Purpose:   ca.Purpose,
				CreatedAt: ca.CreatedAt.Format("2006-01-02"),
			})
			totalCaDeduction += ca.Amount
		}
	}

	// 4. Build empty salary component lines.
	// In the new flow, the frontend handles dynamic addition of components from the master table.
	var components []domain.SalaryComponentLine
	totalAllowance := 0
	totalDeduction := 0

	// 5. Calculate net salary
	netSalary := emp.BaseSalary + totalOvertimePay + totalAllowance - totalDeduction - totalCaDeduction

	return &domain.PayrollSummary{
		EmployeeID:                emp.ID,
		EmployeeName:              emp.Name,
		BaseSalary:                emp.BaseSalary,
		OvertimeInfo:              overtimeInfos,
		OvertimePay:               totalOvertimePay,
		CashAdvances:              caDeductions,
		TotalDeductionCashAdvance: totalCaDeduction,
		SalaryComponents:          components,
		TotalAllowance:            totalAllowance,
		TotalDeduction:            totalDeduction,
		NetSalary:                 netSalary,
	}, nil
}

// calculateHours parses time strings "HH:mm" or "HH:mm:ss" and returns the difference in hours
func calculateHours(startTime, endTime string) float64 {
	layouts := []string{"15:04:05", "15:04"}
	var start, end time.Time
	var err error

	for _, layout := range layouts {
		start, err = time.Parse(layout, startTime)
		if err == nil {
			break
		}
	}
	for _, layout := range layouts {
		end, err = time.Parse(layout, endTime)
		if err == nil {
			break
		}
	}

	diff := end.Sub(start).Hours()
	if diff < 0 {
		diff += 24 // overnight shift
	}
	return math.Round(diff*100) / 100
}
