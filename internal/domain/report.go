package domain

import (
	"context"
	"time"
)

type MonthlySummaryReport struct {
	TotalAttendances int
	TotalOvertimes   int
	TotalLeaves      int
	TotalSickDays    int
	TotalDebts       int
}

type PendingApprovalsResponse struct {
	PendingLeaves       int
	PendingOvertimes    int
	PendingCashAdvances int
}

type DashboardActivity struct {
	ID           int
	EmployeeName string
	Type         string // "attendance", "leave", "overtime", "cash_advance"
	Action       string
	Date         time.Time
	Status       string
}

type CommerceDashboardStats struct {
	TotalRevenue      int
	TotalOrders       int
	PendingOrders     int
	CompletedOrders   int
	TotalShippingCost int
	TotalDiscount     int
	RevenueGraph      []GraphDataPoint
}

type GraphDataPoint struct {
	Label string
	Value int
}

// ReportRepository represents the report's data-access contract
type ReportRepository interface {
	GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*MonthlySummaryReport, error)
	GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error)
	GetPendingApprovals(ctx context.Context) (*PendingApprovalsResponse, error)
	GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]DashboardActivity, int64, error)
	GetCommerceStats(ctx context.Context, filterType string, month, year int) (*CommerceDashboardStats, error)
}

// ReportUsecase represents the report's business logic contract
type ReportUsecase interface {
	GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*MonthlySummaryReport, error)
	GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error)
	GetPendingApprovals(ctx context.Context) (*PendingApprovalsResponse, error)
	GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]DashboardActivity, int64, error)
	GetCommerceStats(ctx context.Context, filterType string, month, year int) (*CommerceDashboardStats, error)
}
