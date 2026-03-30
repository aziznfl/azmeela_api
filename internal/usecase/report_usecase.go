package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type reportUsecase struct {
	reportRepo domain.ReportRepository
	redisRepo  domain.RedisRepository
}

// NewReportUsecase will create a new reportUsecase object representation of domain.ReportUsecase interface
func NewReportUsecase(repo domain.ReportRepository, redisRepo domain.RedisRepository) domain.ReportUsecase {
	return &reportUsecase{
		reportRepo: repo,
		redisRepo:  redisRepo,
	}
}

func (u *reportUsecase) GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*domain.MonthlySummaryReport, error) {
	return u.reportRepo.GetMonthlySummary(ctx, employeeID, month, year)
}

func (u *reportUsecase) GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error) {
	cacheKey := "dashboard_stats:all"
	if employeeID != nil {
		cacheKey = fmt.Sprintf("dashboard_stats:%d", *employeeID)
	}

	var stats map[string]interface{}
	err := u.redisRepo.Get(ctx, cacheKey, &stats)
	if err == nil {
		return stats, nil
	}

	stats, err = u.reportRepo.GetDashboardStats(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	_ = u.redisRepo.Set(ctx, cacheKey, stats, 5*time.Minute)

	return stats, nil
}

func (u *reportUsecase) GetPendingApprovals(ctx context.Context) (*domain.PendingApprovalsResponse, error) {
	return u.reportRepo.GetPendingApprovals(ctx)
}

func (u *reportUsecase) GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]domain.DashboardActivity, int64, error) {
	return u.reportRepo.GetRecentActivities(ctx, employeeID, page, pageSize)
}

func (u *reportUsecase) GetCommerceStats(ctx context.Context, filterType string, month, year int) (*domain.CommerceDashboardStats, error) {
	cacheKey := fmt.Sprintf("commerce_stats:%s:%d:%d", filterType, month, year)

	var stats domain.CommerceDashboardStats
	err := u.redisRepo.Get(ctx, cacheKey, &stats)
	if err == nil {
		return &stats, nil
	}

	result, err := u.reportRepo.GetCommerceStats(ctx, filterType, month, year)
	if err != nil {
		return nil, err
	}

	// Commerce stats can be cached longer, e.g., 10 minutes
	_ = u.redisRepo.Set(ctx, cacheKey, result, 10*time.Minute)

	return result, nil
}
