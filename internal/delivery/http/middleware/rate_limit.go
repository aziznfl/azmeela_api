package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter returns a gin middleware that limits requests
// to 25 requests per 5 seconds. (5 requests per second, burst 25)
func RateLimiter() gin.HandlerFunc {
	// r = rate, b = burst
	// 5 = requests per second, which means 25 requests per 5 seconds
	limiter := rate.NewLimiter(5, 25)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
