package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter returns a gin middleware that limits requests.
// It supports a default limit and a special limit for specific paths.
func RateLimiter(specialPaths ...string) gin.HandlerFunc {
	defaultLimiter := rate.NewLimiter(rate.Limit(5), 20)
	specialLimiter := rate.NewLimiter(rate.Limit(5), 100)

	return func(c *gin.Context) {
		limiter := defaultLimiter

		// Check if the current path is in the special paths list
		for _, path := range specialPaths {
			if c.FullPath() == path {
				limiter = specialLimiter
				break
			}
		}

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
