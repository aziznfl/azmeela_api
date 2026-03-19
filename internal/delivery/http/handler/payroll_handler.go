package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	Usecase    domain.PayrollUsecase
	EmployeeUC domain.EmployeeUsecase
}

func NewPayrollHandler(us domain.PayrollUsecase, empUC domain.EmployeeUsecase) *PayrollHandler {
	return &PayrollHandler{
		Usecase:    us,
		EmployeeUC: empUC,
	}
}

// Generate godoc
// @Summary      Generate payroll for all employees
// @Description  Calculate salary for all employees for a given month/year (superadmin only)
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        month  query  int  true  "Month (1-12)"
// @Param        year   query  int  true  "Year"
// @Success      200    {array}  domain.PayrollSummary
// @Router       /payroll [get]
func (h *PayrollHandler) Generate(c *gin.Context) {
	ctx := c.Request.Context()

	// Only superadmin can generate payroll
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUC.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengidentifikasi peran pengguna")
		return
	}

	if employee.TypeID != 1 {
		ErrorResponse(c, http.StatusForbidden, "Hanya Superadmin yang dapat membuat slip gaji")
		return
	}

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		ErrorResponse(c, http.StatusBadRequest, "Parameter bulan tidak valid (1-12)")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	results, err := h.Usecase.Generate(ctx, month, year)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data payroll berhasil dibuat", dto.ToPayrollListResponse(results), gin.H{
		"month": month,
		"year":  year,
	})
}

// GenerateByEmployee godoc
// @Summary      Generate payroll for a specific employee
// @Description  Calculate salary for one employee for a given month/year (superadmin only)
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        id     path   int  true  "Employee ID"
// @Param        month  query  int  true  "Month (1-12)"
// @Param        year   query  int  true  "Year"
// @Success      200    {object}  domain.PayrollSummary
// @Router       /payroll/{id} [get]
func (h *PayrollHandler) GenerateByEmployee(c *gin.Context) {
	ctx := c.Request.Context()

	// Only superadmin
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	user, err := h.EmployeeUC.GetByID(ctx, payload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to identify user role"})
		return
	}

	if user.TypeID != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can generate payroll"})
		return
	}

	empID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month parameter (1-12)"})
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	result, err := h.Usecase.GenerateByEmployee(ctx, empID, month, year)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data payroll karyawan berhasil dibuat", dto.ToPayrollSummaryResponse(result), gin.H{
		"month": month,
		"year":  year,
	})
}
