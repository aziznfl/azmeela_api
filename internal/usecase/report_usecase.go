package usecase

import (
	"context"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type reportUsecase struct {
	reportRepo domain.ReportRepository
}

// NewReportUsecase will create a new reportUsecase object representation of domain.ReportUsecase interface
func NewReportUsecase(repo domain.ReportRepository) domain.ReportUsecase {
	return &reportUsecase{
		reportRepo: repo,
	}
}

func (u *reportUsecase) GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*domain.MonthlySummaryReport, error) {
	return u.reportRepo.GetMonthlySummary(ctx, employeeID, month, year)
}

func (u *reportUsecase) GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error) {
	return u.reportRepo.GetDashboardStats(ctx, employeeID)
}

func (u *reportUsecase) GetPendingApprovals(ctx context.Context) (*domain.PendingApprovalsResponse, error) {
	return u.reportRepo.GetPendingApprovals(ctx)
}

func (u *reportUsecase) GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]domain.DashboardActivity, int64, error) {
	return u.reportRepo.GetRecentActivities(ctx, employeeID, page, pageSize)
}

func (u *reportUsecase) GetCommerceStats(ctx context.Context, filterType string, month, year int) (*domain.CommerceDashboardStats, error) {
	return u.reportRepo.GetCommerceStats(ctx, filterType, month, year)
}
