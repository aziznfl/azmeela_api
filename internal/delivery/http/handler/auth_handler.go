package handler

import (
	"net/http"
	"strings"

	"github.com/azmeela/sispeg-api/internal/delivery/http/dto"
	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase domain.AuthUsecase
}

// NewAuthHandler will initialize the auth handler
func NewAuthHandler(us domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		Usecase: us,
	}
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	// Set access token cookie
	c.SetCookie(
		"access_token",
		accessToken,
		3600*24, // 24 hours or match your JWT duration
		"/",
		"",
		false, // Set to true in production (Secure)
		true,  // HttpOnly
	)

	// Set refresh token cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		3600*24*7, // 7 days
		"/",
		"",
		false, // Set to true in production (Secure)
		true,  // HttpOnly
	)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}

// Login godoc
// @Summary      Login user
// @Description  Logs in a user using username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.LoginRequest  true  "Login Credentials"
// @Success      200          {object}  domain.AuthResponse
// @Failure      401          {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "")
		return
	}

	ctx := c.Request.Context()
	// Map to domain for usecase
	domainReq := domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := h.Usecase.Login(ctx, &domainReq)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Username atau password yang Anda masukkan salah.")
		return
	}

	h.setAuthCookies(c, res.AccessToken, res.RefreshToken)
	SuccessResponse(c, http.StatusOK, "Login berhasil", dto.ToAuthResponse(res.AccessToken, res.RefreshToken, &res.User))
}

// Refresh godoc
// @Summary      Refresh token
// @Description  Get a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      domain.RefreshRequest  true  "Refresh Token Request"
// @Success      200            {object}  domain.AuthResponse
// @Failure      401            {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		// Fallback to cookie
		cookieToken, err := c.Cookie("refresh_token")
		if err == nil {
			req.RefreshToken = cookieToken
		}
	}

	if req.RefreshToken == "" {
		ErrorResponse(c, http.StatusUnauthorized, "Refresh token is required")
		return
	}

	ctx := c.Request.Context()
	res, err := h.Usecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "")
		return
	}

	h.setAuthCookies(c, res.AccessToken, res.RefreshToken)
	SuccessResponse(c, http.StatusOK, "Token berhasil diperbarui", dto.ToAuthResponse(res.AccessToken, res.RefreshToken, &res.User))
}

// Logout godoc
// @Summary      Logout user
// @Description  Logs out a user and revokes their refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      domain.RefreshRequest  true  "Refresh Token Request"
// @Success      200            {object}  map[string]interface{}
// @Failure      401            {object}  map[string]interface{}
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest

	// try to get from body
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		// fallback to cookie
		cookieToken, err := c.Cookie("refresh_token")
		if err == nil {
			req.RefreshToken = cookieToken
		}

		// fallback to bearer auth headers if cookie misses it
		if req.RefreshToken == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				req.RefreshToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}
	}

	if req.RefreshToken == "" {
		// Even if token is missing, we clear cookies and succeed as it means they are already "out"
		h.clearAuthCookies(c)
		SuccessResponse(c, http.StatusOK, "Berhasil keluar dari sistem", nil)
		return
	}

	ctx := c.Request.Context()
	err := h.Usecase.Logout(ctx, req.RefreshToken)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "")
		return
	}

	h.clearAuthCookies(c)
	SuccessResponse(c, http.StatusOK, "Berhasil keluar dari sistem", nil)
}

// ClearCache godoc
// @Summary      Clear all cache
// @Description  Clears all keys in Redis (Flushes the database)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /cache/clear [get]
func (h *AuthHandler) ClearCache(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.Usecase.ClearCache(ctx)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal membersihkan cache")
		return
	}

	SuccessResponse(c, http.StatusOK, "Berhasil membersihkan seluruh cache", nil)
}
