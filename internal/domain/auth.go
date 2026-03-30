package domain

import (
	"context"
	"time"
)

type LoginRequest struct {
	Username string
	Password string
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
	User         Employee
}

type RefreshRequest struct {
	RefreshToken string
}

// AuthUsecase represent the auth's usecases
type AuthUsecase interface {
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ClearCache(ctx context.Context) error
}

// RedisRepository represent the redis token repository contract
type RedisRepository interface {
	StoreRefreshToken(ctx context.Context, userID int, token string, duration time.Duration) error
	GetRefreshToken(ctx context.Context, token string) (int, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	
	// General Caching
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, target interface{}) error
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
}
