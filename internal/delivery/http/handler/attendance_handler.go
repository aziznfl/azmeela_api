package handler

import (
	"net/http"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Usecase         domain.AttendanceUsecase
	EmployeeUsecase domain.EmployeeUsecase
}

// NewAttendanceHandler will initialize the attendance handler
func NewAttendanceHandler(us domain.AttendanceUsecase, eus domain.EmployeeUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		Usecase:         us,
		EmployeeUsecase: eus,
	}
}

// Fetch godoc
// @Summary      Get attendances
// @Description  Get attendances with optional filtering
// @Tags         attendances
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Attendance
// @Router       /attendances [get]
func (h *AttendanceHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	employee, err := h.EmployeeUsecase.GetByID(ctx, payload.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data karyawan")
		return
	}

	// Example filter parsing
	filter := make(map[string]interface{})

	if employee.TypeID == 1 {
		// Superadmin can query specific employee
		if id := c.Query("employee_id"); id != "" {
			filter["employee_id"] = id
		}
	} else {
		// Normal users can only fetch their own
		filter["employee_id"] = payload.UserID
	}

	if date := c.Query("date"); date != "" {
		filter["presence_date"] = date
	}

	attendances, err := h.Usecase.Fetch(ctx, filter)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar kehadiran berhasil diambil", dto.ToAttendanceListResponseFromDomain(attendances))
}

// GetToday godoc
// @Summary      Get today's attendances
// @Description  Get all attendance records for today
// @Tags         attendances
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Attendance
// @Router       /attendances/today [get]
func (h *AttendanceHandler) GetToday(c *gin.Context) {
	ctx := c.Request.Context()

	attendances, err := h.Usecase.GetTodayAttendances(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	SuccessResponse(c, http.StatusOK, "Daftar kehadiran hari ini berhasil diambil", dto.ToAttendanceListResponseFromDomain(attendances))
}

// ClockIn godoc
// @Summary      Clock-in attendances
// @Description  Create an attendance record for the logged in user
// @Tags         attendances
// @Accept       json
// @Produce      json
// @Success      200  {object}   domain.Attendance
// @Router       /attendances/clock-in [post]
func (h *AttendanceHandler) ClockIn(c *gin.Context) {
	ctx := c.Request.Context()

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	var req dto.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Map DTO to Domain if separate, but here we can just pass the address of mapped domain
	domainReq := domain.AttendanceRequest{
		Location: req.Location,
		Note:     req.Note,
	}

	result, err := h.Usecase.ClockIn(ctx, payload.UserID, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Berhasil masuk (Clock-In)", dto.ToAttendanceResponse(result))
}

// ClockOut godoc
// @Summary      Clock-out attendances
// @Description  Update check-out for the logged in user
// @Tags         attendances
// @Accept       json
// @Produce      json
// @Success      200  {object}   domain.Attendance
// @Router       /attendances/clock-out [post]
func (h *AttendanceHandler) ClockOut(c *gin.Context) {
	ctx := c.Request.Context()

	payloadRaw, _ := c.Get("authorization_payload")
	payload := payloadRaw.(*token.Payload)

	var req dto.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	domainReq := domain.AttendanceRequest{
		Location: req.Location,
		Note:     req.Note,
	}

	result, err := h.Usecase.ClockOut(ctx, payload.UserID, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Berhasil keluar (Clock-Out)", dto.ToAttendanceResponse(result))
}
