package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type SalaryVariableHandler struct {
	Usecase domain.SalaryVariableUsecase
}

func NewSalaryVariableHandler(u domain.SalaryVariableUsecase) *SalaryVariableHandler {
	return &SalaryVariableHandler{
		Usecase: u,
	}
}

// Fetch godoc
// @Summary      Get salary variables
// @Description  Get a list of all salary variables (allowances & deductions)
// @Tags         salary-variables
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.SalaryVariable
// @Router       /salary-variables [get]
func (h *SalaryVariableHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	variables, err := h.Usecase.Fetch(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar variabel gaji berhasil diambil", dto.ToSalaryVariableListResponse(variables))
}

// GetByID godoc
// @Summary      Get salary variable by ID
// @Description  Get a single salary variable by its ID
// @Tags         salary-variables
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Salary Variable ID"
// @Success      200  {object}  domain.SalaryVariable
// @Router       /salary-variables/{id} [get]
func (h *SalaryVariableHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	sv, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Variabel gaji tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data variabel gaji berhasil diambil", dto.ToSalaryVariableResponse(sv))
}

// Store godoc
// @Summary      Create salary variable
// @Description  Create a new salary variable (allowance or deduction)
// @Tags         salary-variables
// @Accept       json
// @Produce      json
// @Param        request  body      domain.SalaryVariableRequest  true  "Salary variable request"
// @Success      201      {object}  domain.SalaryVariable
// @Router       /salary-variables [post]
func (h *SalaryVariableHandler) Store(c *gin.Context) {
	var req dto.SalaryVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	domainReq := domain.SalaryVariableRequest{
		Name:  req.Name,
		Type:  req.Type,
		Value: req.Value,
	}
	result, err := h.Usecase.Create(ctx, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Variabel gaji berhasil dibuat", dto.ToSalaryVariableResponse(result))
}

// Update godoc
// @Summary      Update salary variable
// @Description  Update an existing salary variable
// @Tags         salary-variables
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Salary Variable ID"
// @Param        request  body      domain.SalaryVariableRequest  true  "Salary variable request"
// @Success      200      {object}  domain.SalaryVariable
// @Router       /salary-variables/{id} [put]
func (h *SalaryVariableHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.SalaryVariableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	domainReq := domain.SalaryVariableRequest{
		Name:  req.Name,
		Type:  req.Type,
		Value: req.Value,
	}
	result, err := h.Usecase.Update(ctx, id, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Variabel gaji berhasil diperbarui", dto.ToSalaryVariableResponse(result))
}

// Delete godoc
// @Summary      Delete salary variable
// @Description  Delete a salary variable by ID
// @Tags         salary-variables
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Salary Variable ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /salary-variables/{id} [delete]
func (h *SalaryVariableHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	SuccessResponse(c, http.StatusOK, "Variabel gaji berhasil dihapus", nil)
}
