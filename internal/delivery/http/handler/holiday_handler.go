package handler

import (
	"net/http"
	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type HolidayHandler struct {
	Usecase domain.HolidayUsecase
}

func NewHolidayHandler(u domain.HolidayUsecase) *HolidayHandler {
	return &HolidayHandler{
		Usecase: u,
	}
}

// Fetch godoc
// @Summary      Get holidays
// @Description  Get a list of holidays
// @Tags         holidays
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Holiday
// @Router       /holidays [get]
func (h *HolidayHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()
	filter := make(map[string]interface{})

	holidays, err := h.Usecase.Fetch(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar hari libur berhasil diambil", dto.ToHolidayListResponse(holidays))
}

// Store godoc
// @Summary      Create holiday
// @Description  Create a new holiday
// @Tags         holidays
// @Accept       json
// @Produce      json
// @Param        request  body      domain.HolidayRequest  true  "Holiday request details"
// @Success      201      {object}  domain.Holiday
// @Router       /holidays [post]
func (h *HolidayHandler) Store(c *gin.Context) {
	var req dto.HolidayRequest
	if !BindJSON(c, &req) {
		return
	}

	ctx := c.Request.Context()
	domainReq := domain.HolidayRequest{
		HolidayDate: req.HolidayDate,
		Description: req.Description,
		IsRecurring: req.IsRecurring,
	}
	result, err := h.Usecase.Create(ctx, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Hari libur berhasil dibuat", dto.ToHolidayResponse(result))
}

// Update godoc
// @Summary      Update holiday
// @Description  Update an existing holiday
// @Tags         holidays
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Holiday ID"
// @Param        request  body      domain.HolidayRequest  true  "Holiday details"
// @Success      200      {object}  domain.Holiday
// @Router       /holidays/{id} [put]
func (h *HolidayHandler) Update(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}

	var req dto.HolidayRequest
	if !BindJSON(c, &req) {
		return
	}

	ctx := c.Request.Context()
	domainReq := domain.HolidayRequest{
		HolidayDate: req.HolidayDate,
		Description: req.Description,
		IsRecurring: req.IsRecurring,
	}
	result, err := h.Usecase.Update(ctx, id, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Hari libur berhasil diperbarui", dto.ToHolidayResponse(result))
}

// Delete godoc
// @Summary      Delete holiday
// @Description  Delete a holiday
// @Tags         holidays
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Holiday ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /holidays/{id} [delete]
func (h *HolidayHandler) Delete(c *gin.Context) {
	id, ok := ParseID(c, "id")
	if !ok {
		return
	}

	ctx := c.Request.Context()
	err := h.Usecase.Delete(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Hari libur berhasil dihapus", nil)
}
