package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Usecase domain.ProductUsecase
}

func NewProductHandler(u domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{Usecase: u}
}

func (h *ProductHandler) GetInventory(c *gin.Context) {
	ctx := c.Request.Context()

	filter := make(map[string]interface{})
	if productCodeID := c.Query("product_code_id"); productCodeID != "" {
		id, _ := strconv.Atoi(productCodeID)
		filter["product_code_id"] = id
	}
	if typeID := c.Query("product_type_id"); typeID != "" {
		id, _ := strconv.Atoi(typeID)
		filter["product_type_id"] = id
	}
	if customerTypeID := c.Query("customer_type_id"); customerTypeID != "" {
		id, _ := strconv.Atoi(customerTypeID)
		filter["customer_type_id"] = id
	}

	prices, err := h.Usecase.GetInventoryList(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data inventori berhasil diambil", dto.ToProductCodeListResponse(prices))
}

func (h *ProductHandler) GetCodes(c *gin.Context) {
	ctx := c.Request.Context()

	filter := make(map[string]interface{})
	if typeID := c.Query("product_type_id"); typeID != "" {
		id, _ := strconv.Atoi(typeID)
		filter["product_type_id"] = id
	}

	codes, err := h.Usecase.GetInventoryList(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar kode produk berhasil diambil", dto.ToProductCodeListResponse(codes))
}

func (h *ProductHandler) GetCodesWithTypes(c *gin.Context) {
	ctx := c.Request.Context()

	filter := make(map[string]interface{})
	if customerTypeID := c.Query("customer_type_id"); customerTypeID != "" {
		id, _ := strconv.Atoi(customerTypeID)
		filter["customer_type_id"] = id
	}

	codes, err := h.Usecase.GetCodesWithTypes(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar kode produk dengan tipe berhasil diambil", codes)
}

func (h *ProductHandler) GetTypes(c *gin.Context) {
	ctx := c.Request.Context()
	types, err := h.Usecase.GetProductTypes(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar tipe produk berhasil diambil", dto.ToProductTypeListResponse(types))
}

func (h *ProductHandler) GetSizes(c *gin.Context) {
	ctx := c.Request.Context()
	sizes, err := h.Usecase.GetProductSizes(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar ukuran produk berhasil diambil", dto.ToProductSizeListResponse(sizes))
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	payload := c.MustGet("authorization_payload").(*token.Payload)

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Usecase.UpdateStock(ctx, id, req.Quantity, payload.UserID); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Stok berhasil diperbarui", nil)
}

func (h *ProductHandler) CreateType(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.ProductType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.CreateProductType(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe produk berhasil dibuat", dto.ToProductTypeResponse(&req))
}

func (h *ProductHandler) UpdateType(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	var req domain.ProductType
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.Usecase.UpdateProductType(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe produk berhasil diperbarui", dto.ToProductTypeResponse(&req))
}

func (h *ProductHandler) DeleteType(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Usecase.DeleteProductType(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe produk berhasil dihapus", nil)
}

func (h *ProductHandler) CreateCode(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.ProductCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.CreateProductCode(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Kode produk berhasil dibuat", dto.ToProductCodeResponse(&req))
}

func (h *ProductHandler) UpdateCode(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	var req domain.ProductCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.Usecase.UpdateProductCode(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Kode produk berhasil diperbarui", dto.ToProductCodeResponse(&req))
}

func (h *ProductHandler) DeleteCode(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Usecase.DeleteProductCode(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Kode produk berhasil dihapus", nil)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.CreateProduct(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Produk berhasil dibuat", dto.ToProductResponse(&req))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	var req domain.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.Usecase.UpdateProduct(ctx, &req); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Produk berhasil diperbarui", dto.ToProductResponse(&req))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Usecase.DeleteProduct(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Produk berhasil dihapus", nil)
}

func (h *ProductHandler) GetStockLogs(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	logs, err := h.Usecase.GetStockLogs(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Log stok berhasil diambil", dto.ToProductStockLogListResponse(logs))
}

func (h *ProductHandler) GetColors(c *gin.Context) {
	ctx := c.Request.Context()
	productCodeID, _ := strconv.Atoi(c.Query("product_code_id"))

	colors, err := h.Usecase.GetProductColors(ctx, productCodeID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar warna produk berhasil diambil", colors)
}

func (h *ProductHandler) GetSizesType(c *gin.Context) {
	ctx := c.Request.Context()
	productID, _ := strconv.Atoi(c.Query("product_id"))
	customerTypeID, _ := strconv.Atoi(c.Query("customer_type_id"))

	sizes, err := h.Usecase.GetProductSizesType(ctx, productID, customerTypeID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar ukuran dan harga produk berhasil diambil", sizes)
}
