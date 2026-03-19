package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a middleware for authenticating requests
func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		accessToken := ""

		if len(authorizationHeader) != 0 {
			fields := strings.Fields(authorizationHeader)
			if len(fields) >= 2 && strings.ToLower(fields[0]) == authorizationTypeBearer {
				accessToken = fields[1]
			}
		}

		// Fallback to cookie if header is empty or invalid
		if accessToken == "" {
			cookieToken, err := ctx.Cookie("access_token")
			if err == nil {
				accessToken = cookieToken
			}
		}

		if accessToken == "" {
			err := errors.New("authentication token is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
