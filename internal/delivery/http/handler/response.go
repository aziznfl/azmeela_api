package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse represents a structured error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	// Map common technical errors to human-readable ones if needed
	humanMessage := message
	switch statusCode {
	case http.StatusInternalServerError:
		humanMessage = "Maaf, terjadi kendala pada sistem kami. Silakan coba beberapa saat lagi."
	case http.StatusUnauthorized:
		humanMessage = "Sesi Anda telah berakhir atau Anda tidak memiliki akses. Silakan login kembali."
	case http.StatusForbidden:
		humanMessage = "Anda tidak memiliki izin untuk melakukan aksi ini."
	case http.StatusNotFound:
		humanMessage = "Data yang Anda cari tidak ditemukan."
	case http.StatusBadRequest:
		if message == "" {
			humanMessage = "Data yang Anda kirimkan tidak valid. Periksa kembali inputan Anda."
		}
	}

	c.JSON(statusCode, Response{
		Status:  "error",
		Message: humanMessage,
	})
	// Stop execution
	c.Abort()
}

// SuccessResponse represents a structured success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta ...interface{}) {
	var m interface{}
	if len(meta) > 0 {
		m = meta[0]
	}

	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    m,
	})
}
