package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Usecase domain.TransactionUsecase
}

func NewTransactionHandler(u domain.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		Usecase: u,
	}
}

func (h *TransactionHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()
	filter := make(map[string]interface{})

	// Filters
	if statusID := c.Query("status_id"); statusID != "" {
		filter["status_id"] = statusID
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		filter["customer_id"] = customerID
	}
	if search := c.Query("search"); search != "" {
		filter["search"] = search
	}

	// Pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transactions, meta, err := h.Usecase.Fetch(ctx, filter, page, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", dto.ToTransactionListResponse(transactions), meta)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	transaction, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Transaksi tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data transaksi berhasil diambil", dto.ToTransactionResponse(transaction))
}

func (h *TransactionHandler) Store(c *gin.Context) {
	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.Create(ctx, req.ToDomain())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Transaksi berhasil dibuat", dto.ToTransactionResponse(result))
}

func (h *TransactionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.Update(ctx, id, req.ToDomain())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Transaksi berhasil diperbarui", dto.ToTransactionResponse(result))
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.Usecase.Delete(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Transaksi berhasil dihapus", nil)
}

func (h *TransactionHandler) GetStatuses(c *gin.Context) {
	ctx := c.Request.Context()
	statuses, err := h.Usecase.GetStatuses(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar status transaksi berhasil diambil", dto.ToTransactionStatusListResponse(statuses))
}

func (h *TransactionHandler) GetLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	logs, err := h.Usecase.GetLogs(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Log transaksi berhasil diambil", dto.ToTransactionLogListResponse(logs))
}

func (h *TransactionHandler) GenerateCode(c *gin.Context) {
	ctx := c.Request.Context()
	code, err := h.Usecase.GenerateTransactionCode(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Kode transaksi berhasil dibuat", gin.H{"code": code})
}
