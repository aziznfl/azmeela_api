package middleware

import (
	"fmt"
	"net/http"

	"github.com/azmeela/sispeg-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware catches panics and returns a friendly error message
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the actual error with stack trace (internal only)
				logger.Log.Error(fmt.Sprintf("PANIC RECOVERED: %v", err))

				// Return a very clean, friendly message to the user
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Maaf, terjadi kesalahan internal pada sistem. Kami telah mencatat kejadian ini.",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
