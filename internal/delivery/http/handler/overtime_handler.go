package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type OvertimeHandler struct {
	Usecase         domain.OvertimeUsecase
	EmployeeUsecase domain.EmployeeUsecase
}

func NewOvertimeHandler(us domain.OvertimeUsecase, eu domain.EmployeeUsecase) *OvertimeHandler {
	return &OvertimeHandler{
		Usecase:         us,
		EmployeeUsecase: eu,
	}
}

// Fetch godoc
// @Summary      Get overtimes
// @Description  Get a list of overtime requests
// @Tags         overtimes
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.OvertimeResponse
// @Router       /overtimes [get]
func (h *OvertimeHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	filter := make(map[string]interface{})
	
	if last7Days := c.Query("last_7_days"); last7Days == "true" {
		filter["last_7_days"] = true
	} else {
		month := c.Query("month")
		year := c.Query("year")
		if month != "" && year != "" {
			if len(month) == 1 {
				month = "0" + month
			}
			filter["month_year"] = year + "-" + month + "-%"
		}
	}

	overtimes, err := h.Usecase.Fetch(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar pengajuan lembur berhasil diambil", dto.ToOvertimeListResponseFromDomain(overtimes))
}

// Store godoc
// @Summary      Request overtime
// @Description  Request a new overtime
// @Tags         overtimes
// @Accept       json
// @Produce      json
// @Param        request  body      domain.OvertimeRequest  true  "Overtime request details"
// @Success      201      {object}  domain.Overtime
// @Router       /overtimes [post]
func (h *OvertimeHandler) Store(c *gin.Context) {
	var req dto.OvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	domainReq := domain.OvertimeRequest{
		Date:        req.Date,
		TimeIn:      req.TimeIn,
		TimeOut:     req.TimeOut,
		Description: req.Description,
	}

	result, err := h.Usecase.RequestOvertime(ctx, payload.UserID, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pengajuan lembur berhasil dikirim", dto.ToOvertimeResponse(result))
}

// UpdateStatus godoc
// @Summary      Update overtime status
// @Description  Approve or reject an overtime request
// @Tags         overtimes
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "Overtime ID"
// @Param        request  body      domain.OvertimeStatusUpdate  true  "Status update payload"
// @Success      200      {object}  map[string]interface{}
// @Router       /overtimes/{id}/status [put]
func (h *OvertimeHandler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	var req dto.OvertimeStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	domainReq := domain.OvertimeStatusUpdate{
		Status: req.Status,
	}

	ctx := c.Request.Context()
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data karyawan")
		return
	}

	if employee.TypeID != 1 {
		ErrorResponse(c, http.StatusForbidden, "Hanya Superadmin yang dapat memperbarui status lembur")
		return
	}

	err = h.Usecase.UpdateStatus(ctx, id, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Pengajuan lembur tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Status pengajuan lembur berhasil diperbarui", nil)
}
