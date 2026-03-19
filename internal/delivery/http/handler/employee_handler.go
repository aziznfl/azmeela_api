package handler

import (
	"net/http"
	"strconv"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	Usecase domain.EmployeeUsecase
}

// NewEmployeeHandler will initialize the employees handler
func NewEmployeeHandler(us domain.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		Usecase: us,
	}
}

// Fetch godoc
// @Summary      Get all employees
// @Description  Get a list of all employees
// @Tags         employees
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Employee
// @Router       /employees [get]
func (h *EmployeeHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	employees, err := h.Usecase.Fetch(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar karyawan berhasil diambil", dto.ToEmployeeListResponse(employees))
}

// GetByID godoc
// @Summary      Get a single employee
// @Description  Get an employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Employee ID"
// @Success      200  {object}  domain.Employee
// @Router       /employees/{id} [get]
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	ctx := c.Request.Context()
	employee, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Karyawan tidak ditemukan")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data karyawan berhasil diambil", dto.ToEmployeeResponse(employee))
}

// Store godoc
// @Summary      Create an employee
// @Description  Create a new employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee  body      domain.Employee  true  "Employee object"
// @Success      201       {object}  domain.Employee
// @Router       /employees [post]
func (h *EmployeeHandler) Store(c *gin.Context) {
	var req dto.EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	emp := req.ToDomain()
	err := h.Usecase.Store(ctx, emp)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusCreated, "Karyawan berhasil ditambahkan", dto.ToEmployeeResponse(emp))
}

// Update godoc
// @Summary      Update an employee
// @Description  Update an employee profile
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id        path      int              true  "Employee ID"
// @Param        employee  body      domain.Employee  true  "Employee object"
// @Success      200       {object}  domain.Employee
// @Router       /employees/{id} [put]
func (h *EmployeeHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	var req dto.EmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	
	emp := req.ToDomain()
	emp.ID = id

	ctx := c.Request.Context()
	err = h.Usecase.Update(ctx, emp)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Data karyawan berhasil diperbarui", dto.ToEmployeeResponse(emp))
}

// Delete godoc
// @Summary      Delete an employee
// @Description  Delete an employee profile
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Employee ID"
// @Success      200       {object}  map[string]interface{}
// @Router       /employees/{id} [delete]
func (h *EmployeeHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.Usecase.Delete(ctx, id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Karyawan berhasil dihapus", nil)
}
