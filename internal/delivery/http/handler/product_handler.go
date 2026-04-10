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
	if typeID := c.Query("product_type_id"); typeID != "" {
		id, _ := strconv.Atoi(typeID)
		filter["product_type_id"] = id
	}

	if codeID := c.Query("product_code_id"); codeID != "" {
		id, _ := strconv.Atoi(codeID)
		filter["product_code_id"] = id
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
	SuccessResponse(c, http.StatusOK, "Daftar kode produk dengan tipe berhasil diambil", dto.ToProductCodeWithTypeListResponse(codes))
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
	var req dto.CreateProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	if err := h.Usecase.CreateProductType(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe produk berhasil dibuat", dto.ToProductTypeResponse(entity))
}

func (h *ProductHandler) UpdateType(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.CreateProductTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	entity.ID = id
	if err := h.Usecase.UpdateProductType(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Tipe produk berhasil diperbarui", dto.ToProductTypeResponse(entity))
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
	var req dto.CreateProductCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	if err := h.Usecase.CreateProductCode(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Kode produk berhasil dibuat", dto.ToProductCodeResponse(entity))
}

func (h *ProductHandler) UpdateCode(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.CreateProductCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	entity.ID = id
	if err := h.Usecase.UpdateProductCode(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Kode produk berhasil diperbarui", dto.ToProductCodeResponse(entity))
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
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	if err := h.Usecase.CreateProduct(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Produk berhasil dibuat", dto.ToProductResponse(entity))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	entity.ID = id
	if err := h.Usecase.UpdateProduct(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}
	SuccessResponse(c, http.StatusOK, "Produk berhasil diperbarui", dto.ToProductResponse(entity))
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

func (h *ProductHandler) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()
	productCodeID, _ := strconv.Atoi(c.Query("product_code_id"))

	colors, err := h.Usecase.GetProductColors(ctx, productCodeID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar warna produk berhasil diambil", dto.ToProductColorListResponse(colors))
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

	SuccessResponse(c, http.StatusOK, "Daftar ukuran dan harga produk berhasil diambil", dto.ToProductSizeTypeListResponse(sizes))
}

func (h *ProductHandler) GetSizeWithPrice(c *gin.Context) {
	ctx := c.Request.Context()
	productID, _ := strconv.Atoi(c.Query("product_id"))
	customerTypeID, _ := strconv.Atoi(c.Query("customer_type_id"))

	sizes, err := h.Usecase.GetProductSizesType(ctx, productID, customerTypeID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar ukuran dan harga produk berhasil diambil", dto.ToProductSizeTypeListResponse(sizes))
}

func (h *ProductHandler) GetAllProductSizes(c *gin.Context) {
	ctx := c.Request.Context()
	sizes, err := h.Usecase.GetAllProductSizes(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar semua ukuran produk berhasil diambil", dto.ToProductSizeListResponse(sizes))
}

func (h *ProductHandler) CreateProductSize(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateProductSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	if err := h.Usecase.CreateProductSize(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Ukuran produk berhasil dibuat", dto.ToProductSizeResponse(entity))
}

func (h *ProductHandler) UpdateProductSize(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	var req dto.CreateProductSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := req.ToEntity()
	entity.ID = id
	if err := h.Usecase.UpdateProductSize(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Ukuran produk berhasil diperbarui", dto.ToProductSizeResponse(entity))
}

func (h *ProductHandler) DeleteProductSize(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Usecase.DeleteProductSize(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Ukuran produk berhasil dihapus", nil)
}

func (h *ProductHandler) CreateProductPrice(c *gin.Context) {
	ctx := c.Request.Context()
	payload := c.MustGet("authorization_payload").(*token.Payload)

	var req dto.CreateProductPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := req.ToEntity()
	entity.AdminID = payload.UserID

	if err := h.Usecase.CreateProductPrice(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Harga produk berhasil dibuat", entity)
}

func (h *ProductHandler) UpdateProductPrice(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	payload := c.MustGet("authorization_payload").(*token.Payload)

	var req dto.CreateProductPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity := req.ToEntity()
	entity.ID = id
	entity.AdminID = payload.UserID

	if err := h.Usecase.UpdateProductPrice(ctx, entity); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Harga produk berhasil diperbarui", entity)
}

func (h *ProductHandler) DeleteProductPrice(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Usecase.DeleteProductPrice(ctx, id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Harga produk berhasil dihapus", nil)
}

