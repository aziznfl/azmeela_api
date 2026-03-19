package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	Usecase domain.LeaveUsecase
}

// NewLeaveHandler will initialize the leave handler
func NewLeaveHandler(us domain.LeaveUsecase) *LeaveHandler {
	return &LeaveHandler{
		Usecase: us,
	}
}

// Fetch godoc
// @Summary      Get leaves
// @Description  Get a list of leave requests
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Param        type   query      string  false  "Type filter (cuti/sakit)"
// @Success      200  {array}   domain.LeaveResponse
// @Router       /leaves [get]
func (h *LeaveHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	filter := make(map[string]interface{})
	if leaveType := c.Query("type"); leaveType != "" {
		if t, err := strconv.Atoi(leaveType); err == nil {
			filter["type"] = t
		}
	}
	
	if upcoming := c.Query("upcoming"); upcoming == "true" {
		filter["upcoming"] = true
	} else if last7Days := c.Query("last_7_days"); last7Days == "true" {
		filter["last_7_days"] = true
	} else {
		// Only apply month/year if not special filters
		month := c.Query("month")
		year := c.Query("year")
		if month != "" && year != "" {
			// PostgreSQL / MySQL generic approach: "YYYY-MM-%"
			// Ensure two digit month
			if len(month) == 1 {
				month = "0" + month
			}
			filter["month_year"] = year + "-" + month + "-%"
		}
	}

	leaves, err := h.Usecase.Fetch(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar pengajuan cuti berhasil diambil", dto.ToLeaveListResponseFromDomain(leaves))
}

// Store godoc
// @Summary      Request Leave
// @Description  Request a new leave or sick day
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Param        request  body      domain.LeaveRequest  true  "Leave request details"
// @Success      201      {object}  domain.Leave
// @Router       /leaves [post]
func (h *LeaveHandler) Store(c *gin.Context) {
	var req dto.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	domainReq := domain.LeaveRequest{
		Type:        req.Type,
		LeaveDate:   req.LeaveDate,
		Durations:   req.Durations,
		Description: req.Description,
	}

	result, err := h.Usecase.RequestLeave(ctx, payload.UserID, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pengajuan cuti berhasil dikirim", dto.ToLeaveResponse(result))
}

// UpdateStatus godoc
// @Summary      Update leave status
// @Description  Approve or reject a leave request
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Leave ID"
// @Param        request  body      domain.LeaveStatusUpdate  true  "Leave status override"
// @Success      200      {object}  map[string]interface{}
// @Router       /leaves/{id}/status [put]
func (h *LeaveHandler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	var req dto.LeaveStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	domainReq := domain.LeaveStatusUpdate{
		Status: req.Status,
	}

	ctx := c.Request.Context()
	err = h.Usecase.UpdateStatus(ctx, id, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Pengajuan cuti tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Status pengajuan cuti berhasil diperbarui", nil)
}
