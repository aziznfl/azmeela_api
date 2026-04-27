package middleware

import (
	"time"

	"github.com/azmeela/sispeg-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware records HTTP request details and logs them using structured logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the log metrics
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 400 {
			logger.Log.Error("HTTP Request",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency.String(),
				"user-agent", c.Request.UserAgent(),
				"error", errorMessage,
			)
		} else {
			logger.Log.Info("HTTP Request",
				"status", statusCode,
				"method", method,
				"path", path,
				"ip", clientIP,
				"latency", latency.String(),
				"user-agent", c.Request.UserAgent(),
			)
		}
	}
}
