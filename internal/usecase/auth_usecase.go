package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/azmeela/sispeg-api/pkg/token"
	"github.com/azmeela/sispeg-api/pkg/utils"
)

type authUsecase struct {
	employeeRepo domain.EmployeeRepository
	redisRepo    domain.RedisRepository
	tokenMaker   token.TokenMaker
}

// NewAuthUsecase will create new an authUsecase object representation of domain.AuthUsecase interface
func NewAuthUsecase(e domain.EmployeeRepository, r domain.RedisRepository, t token.TokenMaker) domain.AuthUsecase {
	return &authUsecase{
		employeeRepo: e,
		redisRepo:    r,
		tokenMaker:   t,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	emp, err := u.employeeRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Codeigniter used MD5 hash for password
	if !utils.VerifyMD5(req.Password, emp.Password) {
		return nil, errors.New("invalid username or password")
	}

	if emp.Active == 0 {
		return nil, errors.New("account is inactive")
	}

	return u.generateAuthTokens(ctx, emp)
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	// Verify refresh token signature
	payload, err := u.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if refresh token exists in redis
	userID, err := u.redisRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil || userID != payload.UserID {
		return nil, errors.New("refresh token already logged out or revoked")
	}

	emp, err := u.employeeRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if emp.Active == 0 {
		return nil, errors.New("account is inactive")
	}

	// Delete old refresh token from redis
	u.redisRepo.DeleteRefreshToken(ctx, refreshToken)

	return u.generateAuthTokens(ctx, emp)
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	// Delete any existing refresh token from Redis
	return u.redisRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (u *authUsecase) generateAuthTokens(ctx context.Context, emp *domain.Employee) (*domain.AuthResponse, error) {
	// durations from rules
	accessTokenDuration := 15 * time.Minute
	refreshTokenDuration := 7 * 24 * time.Hour

	accessToken, err := u.tokenMaker.CreateToken(emp.ID, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenMaker.CreateToken(emp.ID, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	// store refresh token to redis
	err = u.redisRepo.StoreRefreshToken(ctx, emp.ID, refreshToken, refreshTokenDuration)
	if err != nil {
		return nil, errors.New("failed to store token")
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *emp,
	}, nil
}
