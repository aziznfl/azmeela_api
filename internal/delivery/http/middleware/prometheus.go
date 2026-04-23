package middleware

import (
	"strconv"
	"time"

	"github.com/azmeela/sispeg-api/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware is a Gin middleware that records HTTP request metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Use the full path (e.g., /api/v1/employees/:id) instead of the actual URI to avoid cardinality issues
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		metrics.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
