package repository

import (
	"context"
	"sort"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) domain.ReportRepository {
	return &reportRepository{db}
}

// dateRange returns a query filtered by a date column within [start, end).
func (r *reportRepository) dateRange(ctx context.Context, model interface{}, col string, start, end time.Time) *gorm.DB {
	return r.db.WithContext(ctx).Model(model).Where(col+" >= ? AND "+col+" < ?", start, end)
}

// statusLabel maps a numeric status to a human-readable string.
func statusLabel(status int) string {
	switch status {
	case 1:
		return "Disetujui"
	case 2:
		return "Ditolak"
	default:
		return "Pending"
	}
}

// sumColumn returns the SUM of a column using COALESCE.
func sumColumn(q *gorm.DB, col string) int64 {
	var total int64
	q.Select("COALESCE(SUM(" + col + "), 0)").Scan(&total)
	return total
}

// countWhere counts rows with optional extra where clause.
func countWhere(q *gorm.DB, where ...interface{}) int64 {
	var count int64
	if len(where) > 0 {
		q = q.Where(where[0], where[1:]...)
	}
	q.Count(&count)
	return count
}

func (r *reportRepository) GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*domain.MonthlySummaryReport, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	build := func(model interface{}, col string) *gorm.DB {
		q := r.dateRange(ctx, model, col, startDate, endDate)
		if employeeID != nil {
			q = q.Where("id_admin = ?", *employeeID)
		}
		return q
	}

	var summary domain.MonthlySummaryReport
	summary.TotalAttendances = int(countWhere(build(&domain.Attendance{}, "tanggal")))
	summary.TotalOvertimes = int(countWhere(build(&domain.Overtime{}, "tanggal"), "status = 1"))
	summary.TotalLeaves = int(countWhere(build(&domain.Leave{}, "tanggal"), "status = 1 AND grouping = 0"))
	summary.TotalSickDays = int(countWhere(build(&domain.Leave{}, "tanggal"), "status = 1 AND grouping = 1"))
	summary.TotalDebts = int(sumColumn(build(&domain.CashAdvance{}, "tanggal").Where("status = 1"), "jumlah"))

	return &summary, nil
}

func (r *reportRepository) GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error) {
	today := time.Now().UTC().Format("2006-01-02")
	stats := make(map[string]interface{})

	var totalEmployees int64
	r.db.WithContext(ctx).Model(&domain.Employee{}).Where("id_admin_type != 1 AND status_admin = 1").Count(&totalEmployees)
	stats["total_employees"] = totalEmployees

	qPresent := r.db.WithContext(ctx).Model(&domain.Attendance{}).Where("tanggal = ?", today)
	if employeeID != nil {
		qPresent = qPresent.Where("id_admin = ?", *employeeID)
	}
	presentCount := countWhere(qPresent)

	if employeeID == nil {
		if totalEmployees > 0 {
			stats["presence_percentage"] = float64(presentCount) / float64(totalEmployees) * 100
		} else {
			stats["presence_percentage"] = 0
		}
	} else {
		stats["presence_percentage"] = presentCount
	}

	qLeaves := r.db.WithContext(ctx).Model(&domain.Leave{}).Where("tanggal = ? AND status = 1 AND grouping = 0", today)
	if employeeID != nil {
		qLeaves = qLeaves.Where("id_admin = ?", *employeeID)
	}
	stats["on_leave"] = countWhere(qLeaves)

	return stats, nil
}

func (r *reportRepository) GetPendingApprovals(ctx context.Context) (*domain.PendingApprovalsResponse, error) {
	db := r.db.WithContext(ctx)
	var resp domain.PendingApprovalsResponse

	resp.PendingLeaves = int(countWhere(db.Model(&domain.Leave{}), "status = 0"))
	resp.PendingOvertimes = int(countWhere(db.Model(&domain.Overtime{}), "status = 0"))
	resp.PendingCashAdvances = int(countWhere(db.Model(&domain.CashAdvance{}), "status = 0"))

	return &resp, nil
}

func (r *reportRepository) GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]domain.DashboardActivity, int64, error) {
	sevenDaysAgo := time.Now().UTC().AddDate(0, 0, -7)
	maxRecent := 500

	withEmployee := func(q *gorm.DB, table string) *gorm.DB {
		if employeeID != nil {
			q = q.Where(table+".id_admin = ?", *employeeID)
		}
		return q
	}

	var activities []domain.DashboardActivity

	// Attendances
	var attendances []domain.Attendance
	withEmployee(
		r.db.WithContext(ctx).Model(&domain.Attendance{}).
			Select("t_presensi.*, t_admin.nama_admin as employee_name").
			Joins("JOIN t_admin ON t_admin.id_admin = t_presensi.id_admin").
			Where("t_presensi.tanggal >= ?", sevenDaysAgo).
			Order("t_presensi.tanggal DESC, t_presensi.jam_masuk DESC").
			Limit(maxRecent),
		"t_presensi",
	).Find(&attendances)

	parseTime := func(ts string) time.Time {
		t, err := time.Parse("15:04:05", ts)
		if err != nil {
			t, _ = time.Parse("15:04", ts)
		}
		return t
	}

	for _, a := range attendances {
		tIn := parseTime(a.TimeIn)
		activities = append(activities, domain.DashboardActivity{
			ID: a.ID, EmployeeName: a.EmployeeName, Type: "attendance",
			Action: "Absen Masuk",
			Date:   time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), tIn.Hour(), tIn.Minute(), tIn.Second(), 0, time.UTC),
			Status: "Hadir",
		})
		if a.TimeOut != nil && *a.TimeOut != "" {
			tOut := parseTime(*a.TimeOut)
			activities = append(activities, domain.DashboardActivity{
				ID: a.ID, EmployeeName: a.EmployeeName, Type: "attendance",
				Action: "Absen Pulang",
				Date:   time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), tOut.Hour(), tOut.Minute(), tOut.Second(), 0, time.UTC),
				Status: "Hadir",
			})
		}
	}

	// Leaves
	var leaves []domain.Leave
	withEmployee(
		r.db.WithContext(ctx).Model(&domain.Leave{}).
			Select("t_cuti.*, t_admin.nama_admin as employee_name").
			Joins("JOIN t_admin ON t_admin.id_admin = t_cuti.id_admin").
			Where("t_cuti.tanggal >= ?", sevenDaysAgo).
			Order("t_cuti.tanggal DESC").Limit(maxRecent),
		"t_cuti",
	).Find(&leaves)

	for _, l := range leaves {
		typeStr := "Cuti"
		if l.Type == 1 {
			typeStr = "Sakit"
		}
		activities = append(activities, domain.DashboardActivity{
			ID: l.ID, EmployeeName: l.EmployeeName, Type: "leave",
			Action: "Pengajuan " + typeStr,
			Date:   l.LeaveDate.UTC(),
			Status: statusLabel(l.Status),
		})
	}

	// Overtimes
	var overtimes []domain.Overtime
	withEmployee(
		r.db.WithContext(ctx).Model(&domain.Overtime{}).
			Select("t_lembur.*, t_admin.nama_admin as employee_name").
			Joins("JOIN t_admin ON t_admin.id_admin = t_lembur.id_admin").
			Where("t_lembur.tanggal >= ?", sevenDaysAgo).
			Order("t_lembur.tanggal DESC").Limit(maxRecent),
		"t_lembur",
	).Find(&overtimes)

	for _, o := range overtimes {
		activities = append(activities, domain.DashboardActivity{
			ID: o.ID, EmployeeName: o.EmployeeName, Type: "overtime",
			Action: "Pengajuan Lembur",
			Date:   o.Date.UTC(),
			Status: statusLabel(o.Status),
		})
	}

	// Cash Advances
	var cashAdvances []domain.CashAdvance
	withEmployee(
		r.db.WithContext(ctx).Model(&domain.CashAdvance{}).
			Select("t_kasbon.*, t_admin.nama_admin as employee_name").
			Joins("JOIN t_admin ON t_admin.id_admin = t_kasbon.id_admin").
			Where("t_kasbon.tanggal >= ?", sevenDaysAgo).
			Order("t_kasbon.tanggal DESC").Limit(maxRecent),
		"t_kasbon",
	).Find(&cashAdvances)

	for _, ca := range cashAdvances {
		activities = append(activities, domain.DashboardActivity{
			ID: ca.ID, EmployeeName: ca.EmployeeName, Type: "cash_advance",
			Action: "Pengajuan Kasbon",
			Date:   ca.CreatedAt.UTC(),
			Status: statusLabel(ca.Status),
		})
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Date.After(activities[j].Date)
	})

	totalCount := int64(len(activities))
	start := (page - 1) * pageSize
	if start > len(activities) {
		return []domain.DashboardActivity{}, totalCount, nil
	}
	end := start + pageSize
	if end > len(activities) {
		end = len(activities)
	}

	return activities[start:end], totalCount, nil
}

func (r *reportRepository) GetCommerceStats(ctx context.Context, filterType string, month, year int) (*domain.CommerceDashboardStats, error) {
	var stats domain.CommerceDashboardStats
	now := time.Now().UTC()

	getQuery := func() *gorm.DB {
		base := r.db.WithContext(ctx).Model(&domain.Transaction{})
		switch filterType {
		case "monthly":
			start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			return base.Where("tgl_transaksi >= ? AND tgl_transaksi < ?", start, start.AddDate(0, 1, 0))
		case "yearly":
			start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			return base.Where("tgl_transaksi >= ? AND tgl_transaksi < ?", start, start.AddDate(1, 0, 0))
		default:
			return base.Where("tgl_transaksi >= ?", now.AddDate(0, 0, -7))
		}
	}

	totalRevenue := sumColumn(getQuery(), "total")
	totalDiscount := sumColumn(getQuery(), "diskon")

	stats.TotalOrders = int(countWhere(getQuery()))
	stats.PendingOrders = int(countWhere(getQuery(), "transaksi_status = 1"))
	stats.CompletedOrders = int(countWhere(getQuery(), "transaksi_status >= 104"))
	stats.TotalShippingCost = int(sumColumn(getQuery(), "ongkir"))
	stats.TotalDiscount = int(totalDiscount)
	stats.TotalRevenue = int(totalRevenue - totalDiscount) // value = total - discount

	// Graph Data
	revQuery := func(start, end time.Time) int {
		var v int64
		r.db.WithContext(ctx).Table("t_transaksi").
			Select("COALESCE(SUM(total - diskon), 0)").
			Where("tgl_transaksi >= ? AND tgl_transaksi < ?", start, end).
			Scan(&v)
		return int(v)
	}

	var graphPoints []domain.GraphDataPoint
	switch filterType {
	case "yearly":
		for m := 1; m <= 12; m++ {
			start := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: start.Format("Jan"),
				Value: revQuery(start, start.AddDate(0, 1, 0)),
			})
		}
	case "monthly":
		daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for d := 1; d <= daysInMonth; d++ {
			start := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: start.Format("02 Jan"),
				Value: revQuery(start, start.AddDate(0, 0, 1)),
			})
		}
	default:
		for i := 6; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: start.Format("02 Jan"),
				Value: revQuery(start, start.AddDate(0, 0, 1)),
			})
		}
	}

	stats.RevenueGraph = graphPoints
	return &stats, nil
}
