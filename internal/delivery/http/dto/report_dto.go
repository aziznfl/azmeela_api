package dto

import (
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
)

type MonthlySummaryReport struct {
	TotalAttendances int `json:"total_attendances"`
	TotalOvertimes   int `json:"total_overtimes"`
	TotalLeaves      int `json:"total_leaves"`
	TotalSickDays    int `json:"total_sick_days"`
	TotalDebts       int `json:"total_debts"`
}

type PendingApprovalsResponse struct {
	PendingLeaves       int `json:"pending_leaves"`
	PendingOvertimes    int `json:"pending_overtimes"`
	PendingCashAdvances int `json:"pending_cash_advances"`
}

type DashboardActivityResponse struct {
	ID           int       `json:"id"`
	EmployeeName string    `json:"employee_name"`
	Type         string    `json:"type"` // "attendance", "leave", "overtime", "cash_advance"
	Action       string    `json:"action"`
	Date         time.Time `json:"date"`
	Status       string    `json:"status"`
}

type CommerceDashboardStats struct {
	TotalRevenue      int                   `json:"total_revenue"`
	TotalOrders       int                   `json:"total_orders"`
	PendingOrders     int                   `json:"pending_orders"`
	CompletedOrders   int                   `json:"completed_orders"`
	TotalShippingCost int                   `json:"total_shipping_cost"`
	TotalDiscount     int                   `json:"total_discount"`
	RevenueGraph      []GraphDataPointResponse `json:"revenue_graph"`
}

type GraphDataPointResponse struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

func ToDashboardActivityResponse(a *domain.DashboardActivity) *DashboardActivityResponse {
	if a == nil {
		return nil
	}
	return &DashboardActivityResponse{
		ID:           a.ID,
		EmployeeName: a.EmployeeName,
		Type:         a.Type,
		Action:       a.Action,
		Date:         a.Date,
		Status:       a.Status,
	}
}

func ToDashboardActivityListResponse(items []domain.DashboardActivity) []*DashboardActivityResponse {
	resps := make([]*DashboardActivityResponse, len(items))
	for i, item := range items {
		resps[i] = ToDashboardActivityResponse(&item)
	}
	return resps
}

func ToCommerceDashboardStatsResponse(s *domain.CommerceDashboardStats) *CommerceDashboardStats {
	if s == nil {
		return nil
	}
	
	graph := make([]GraphDataPointResponse, len(s.RevenueGraph))
	for i, g := range s.RevenueGraph {
		graph[i] = GraphDataPointResponse(g)
	}

	return &CommerceDashboardStats{
		TotalRevenue:      s.TotalRevenue,
		TotalOrders:       s.TotalOrders,
		PendingOrders:     s.PendingOrders,
		CompletedOrders:   s.CompletedOrders,
		TotalShippingCost: s.TotalShippingCost,
		TotalDiscount:     s.TotalDiscount,
		RevenueGraph:      graph,
	}
}
