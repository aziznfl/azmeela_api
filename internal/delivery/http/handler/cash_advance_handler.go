package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type CashAdvanceHandler struct {
	Usecase    domain.CashAdvanceUsecase
	EmployeeUC domain.EmployeeUsecase
}

func NewCashAdvanceHandler(us domain.CashAdvanceUsecase, empUC domain.EmployeeUsecase) *CashAdvanceHandler {
	return &CashAdvanceHandler{
		Usecase:    us,
		EmployeeUC: empUC,
	}
}

// Fetch godoc
// @Summary      Get cash advances
// @Description  Superadmin (type_id=1) sees all cash advances. Normal admin sees only their own.
// @Tags         cash-advances
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.CashAdvanceResponse
// @Router       /cash-advances [get]
func (h *CashAdvanceHandler) Fetch(c *gin.Context) {
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

	advances, err := h.Usecase.Fetch(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar kasbon berhasil diambil", dto.ToCashAdvanceListResponseFromDomain(advances))
}

// Store godoc
// @Summary      Request cash advance
// @Description  Create a new cash advance request (any authenticated user)
// @Tags         cash-advances
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CashAdvanceRequest  true  "Cash advance request details"
// @Success      201      {object}  domain.CashAdvance
// @Router       /cash-advances [post]
func (h *CashAdvanceHandler) Store(c *gin.Context) {
	var req dto.CashAdvanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	domainReq := domain.CashAdvanceRequest{
		Amount:  req.Amount,
		Purpose: req.Purpose,
	}

	result, err := h.Usecase.RequestCashAdvance(ctx, payload.UserID, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pengajuan kasbon berhasil dikirim", dto.ToCashAdvanceResponse(result))
}

// UpdateStatus godoc
// @Summary      Update cash advance status
// @Description  Approve or reject a cash advance request (superadmin only)
// @Tags         cash-advances
// @Accept       json
// @Produce      json
// @Param        id       path      int                             true  "Cash Advance ID"
// @Param        request  body      domain.CashAdvanceStatusUpdate  true  "Status update payload"
// @Success      200      {object}  map[string]interface{}
// @Router       /cash-advances/{id}/status [put]
func (h *CashAdvanceHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Only superadmin can update status
	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUC.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengidentifikasi peran pengguna")
		return
	}

	if employee.TypeID != 1 {
		ErrorResponse(c, http.StatusForbidden, "Hanya Superadmin yang dapat memperbarui status kasbon")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	var req dto.CashAdvanceStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	domainReq := domain.CashAdvanceStatusUpdate{
		Status: req.Status,
	}

	err = h.Usecase.UpdateStatus(ctx, id, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Data kasbon tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Status kasbon berhasil diperbarui", nil)
}

// AddPayment godoc
// @Summary      Add cash advance payment
// @Description  Record a payment toward a cash advance (debt)
// @Tags         cash-advances
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CashAdvancePayment  true  "Payment details"
// @Success      201      {object}  map[string]interface{}
// @Router       /cash-advances/payment [post]
func (h *CashAdvanceHandler) AddPayment(c *gin.Context) {
	var req dto.CashAdvancePayment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	domainReq := domain.CashAdvancePayment{
		EmployeeID: req.EmployeeID,
		Date:       req.Date,
		Amount:     req.Amount,
	}

	ctx := c.Request.Context()
	err := h.Usecase.AddPayment(ctx, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pembayaran kasbon berhasil dicatat", nil)
}
