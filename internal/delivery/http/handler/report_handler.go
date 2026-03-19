package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Usecase         domain.ReportUsecase
	EmployeeUsecase domain.EmployeeUsecase
}

func NewReportHandler(u domain.ReportUsecase, eu domain.EmployeeUsecase) *ReportHandler {
	return &ReportHandler{
		Usecase:         u,
		EmployeeUsecase: eu,
	}
}

// GetMonthlySummary godoc
// @Summary      Get monthly summary report
// @Description  Get a summary of attendances, overtimes, leaves, and debts for a specific month
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        month  query      int  false  "Month (1-12), defaults to current"
// @Param        year   query      int  false  "Year (YYYY), defaults to current"
// @Success      200    {object}  domain.MonthlySummaryReport
// @Router       /reports/monthly-summary [get]
func (h *ReportHandler) GetMonthlySummary(c *gin.Context) {
	ctx := c.Request.Context()

	now := time.Now()
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data karyawan")
		return
	}

	var pEmployeeID *int
	if employee.TypeID == 1 {
		// If superadmin and query param provided
		if queryID := c.Query("employee_id"); queryID != "" {
			if id, err := strconv.Atoi(queryID); err == nil {
				pEmployeeID = &id
			}
		}
		// If omitted, pEmployeeID remains nil -> all employees stats
	} else {
		// Normal user, only their own
		pEmployeeID = &payload.UserID
	}

	report, err := h.Usecase.GetMonthlySummary(ctx, pEmployeeID, month, year)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Ringkasan bulanan berhasil diambil", dto.MonthlySummaryReport(*report))
}

// GetDashboardStats godoc
// @Summary      Get dashboard stats
// @Description  Get a summary of stats for the dashboard showing all users or specific user
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        employee_id  query      int  false  "Employee ID (Superadmin only)"
// @Success      200    {object}  map[string]interface{}
// @Router       /reports/dashboard-stats [get]
func (h *ReportHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data karyawan")
		return
	}

	var pEmployeeID *int
	if employee.TypeID == 1 {
		// If superadmin and query param provided
		if queryID := c.Query("employee_id"); queryID != "" {
			if id, err := strconv.Atoi(queryID); err == nil {
				pEmployeeID = &id
			}
		}
		// If omitted, pEmployeeID remains nil -> all employees stats
	} else {
		// Normal user, only their own
		pEmployeeID = &payload.UserID
	}

	stats, err := h.Usecase.GetDashboardStats(ctx, pEmployeeID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Statistik dashboard berhasil diambil", stats)
}

// GetPendingApprovals godoc
// @Summary      Get pending approvals
// @Description  Get a count of pending requests for superadmin notifications
// @Tags         reports
// @Accept       json
// @Produce      json
// @Success      200    {object}  domain.PendingApprovalsResponse
// @Router       /reports/pending-approvals [get]
func (h *ReportHandler) GetPendingApprovals(c *gin.Context) {
	ctx := c.Request.Context()

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil || employee.TypeID != 1 {
		ErrorResponse(c, http.StatusForbidden, "Hanya Superadmin yang dapat melihat persetujuan tertunda")
		return
	}

	resp, err := h.Usecase.GetPendingApprovals(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar persetujuan tertunda berhasil diambil", dto.PendingApprovalsResponse(*resp))
}

// GetRecentActivities godoc
// @Summary      Get recent activities
// @Description  Get activities from last 7 days (attendances, leaves, overtimes, cash advances)
// @Tags         reports
// @Accept       json
// @Produce      json
// @Success      200    {array}   domain.DashboardActivity
// @Router       /reports/recent-activities [get]
func (h *ReportHandler) GetRecentActivities(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data karyawan")
		return
	}

	var pEmployeeID *int
	// Superadmin can filter by employee_id if they want
	if employee.TypeID == 1 {
		if queryID := c.Query("employee_id"); queryID != "" {
			if id, err := strconv.Atoi(queryID); err == nil {
				pEmployeeID = &id
			}
		}
	} else {
		// Normal employees see everyone's activities as requested
		pEmployeeID = nil
	}

	activities, total, err := h.Usecase.GetRecentActivities(ctx, pEmployeeID, page, pageSize)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar aktivitas terbaru berhasil diambil", dto.ToDashboardActivityListResponse(activities), gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// GetCommerceStats godoc
// @Summary      Get commerce stats
// @Description  Get a summary of sales, orders, and revenue graph
// @Tags         reports
// @Accept       json
// @Produce      json
// @Success      200    {object}  domain.CommerceDashboardStats
// @Router       /reports/commerce-stats [get]
func (h *ReportHandler) GetCommerceStats(c *gin.Context) {
	ctx := c.Request.Context()

	filterType := c.DefaultQuery("filter_type", "last-7-days")
	now := time.Now()
	
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		month = int(now.Month())
	}
	
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		year = now.Year()
	}

	stats, err := h.Usecase.GetCommerceStats(ctx, filterType, month, year)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Statistik commerce berhasil diambil", dto.ToCommerceDashboardStatsResponse(stats))
}
