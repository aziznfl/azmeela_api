package handler

import (
	"net/http"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type AdminTypeHandler struct {
	Usecase domain.AdminTypeUsecase
}

func NewAdminTypeHandler(u domain.AdminTypeUsecase) *AdminTypeHandler {
	return &AdminTypeHandler{Usecase: u}
}

// Fetch godoc
// @Summary      Get all admin types
// @Description  Get a list of all employee roles/types
// @Tags         admin-types
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.AdminType
// @Router       /admin-types [get]
func (h *AdminTypeHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	types, err := h.Usecase.Fetch(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar tipe admin berhasil diambil", dto.ToAdminTypeListResponse(types))
}
