package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	Usecase domain.CustomerUsecase
}

func NewCustomerHandler(u domain.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		Usecase: u,
	}
}

func (h *CustomerHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()
	filter := make(map[string]interface{})

	// Basic filters
	if typeID := c.Query("type_id"); typeID != "" {
		filter["customer_type_id"] = typeID
	}
	if search := c.Query("search"); search != "" {
		filter["search"] = search
	}

	// Pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	customers, meta, err := h.Usecase.Fetch(ctx, filter, page, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar pelanggan berhasil diambil", dto.ToCustomerListResponse(customers), meta)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	customer, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Pelanggan tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data pelanggan berhasil diambil", dto.ToCustomerResponse(customer))
}

func (h *CustomerHandler) Store(c *gin.Context) {
	var req domain.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.Create(ctx, &req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pelanggan berhasil ditambahkan", dto.ToCustomerResponse(result))
}

func (h *CustomerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req domain.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.Update(ctx, id, &req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data pelanggan berhasil diperbarui", dto.ToCustomerResponse(result))
}

func (h *CustomerHandler) Delete(c *gin.Context) {
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

	SuccessResponse(c, http.StatusOK, "Pelanggan berhasil dihapus", nil)
}

func (h *CustomerHandler) GetTypes(c *gin.Context) {
	ctx := c.Request.Context()
	types, err := h.Usecase.GetTypes(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar tipe pelanggan berhasil diambil", dto.ToCustomerTypeListResponse(types))
}

func (h *CustomerHandler) CreateType(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.CustomerType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.CreateType(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe pelanggan berhasil dibuat", dto.ToCustomerTypeResponse(&req))
}

func (h *CustomerHandler) UpdateType(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.CustomerType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.UpdateType(ctx, id, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe pelanggan berhasil diperbarui", dto.ToCustomerTypeResponse(&req))
}

func (h *CustomerHandler) DeleteType(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Usecase.DeleteType(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe pelanggan berhasil dihapus", nil)
}

func (h *CustomerHandler) GetAddresses(c *gin.Context) {
	customerID, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	addresses, err := h.Usecase.GetAddresses(ctx, customerID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar alamat pelanggan berhasil diambil", dto.ToCustomerAddressListResponse(addresses))
}

func (h *CustomerHandler) CreateAddress(c *gin.Context) {
	var req domain.CustomerAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.CreateAddress(ctx, &req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, "Alamat berhasil ditambahkan", dto.ToCustomerAddressResponse(result))
}

func (h *CustomerHandler) UpdateAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("address_id"))
	var req domain.CustomerAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	result, err := h.Usecase.UpdateAddress(ctx, id, &req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Alamat berhasil diperbarui", dto.ToCustomerAddressResponse(result))
}

func (h *CustomerHandler) DeleteAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("address_id"))
	ctx := c.Request.Context()
	err := h.Usecase.DeleteAddress(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Alamat berhasil dihapus", nil)
}
