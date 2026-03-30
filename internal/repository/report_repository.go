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

// NewReportRepository will create an object that represents the domain.ReportRepository interface
func NewReportRepository(db *gorm.DB) domain.ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetMonthlySummary(ctx context.Context, employeeID *int, month, year int) (*domain.MonthlySummaryReport, error) {
	var summary domain.MonthlySummaryReport

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	
	db := r.db.WithContext(ctx)

	// Senior BE Optimization: Range queries are much faster than LIKE on indexed date columns
	buildQuery := func(model interface{}, dateColumn string) *gorm.DB {
		q := db.Model(model).Where(dateColumn+" >= ? AND "+dateColumn+" < ?", startDate, endDate)
		if employeeID != nil {
			q = q.Where("id_admin = ?", *employeeID)
		}
		return q
	}

	// Count attendances
	var count int64
	buildQuery(&domain.Attendance{}, "tanggal").Count(&count)
	summary.TotalAttendances = int(count)

	// Count overtimes
	buildQuery(&domain.Overtime{}, "tanggal").Where("status = 1").Count(&count)
	summary.TotalOvertimes = int(count)

	// Count leaves (cuti), grouping = 0
	buildQuery(&domain.Leave{}, "tanggal").Where("status = 1 AND grouping = 0").Count(&count)
	summary.TotalLeaves = int(count)

	// Count sick days (sakit), grouping = 1
	var sickCount int64
	buildQuery(&domain.Leave{}, "tanggal").Where("status = 1 AND grouping = 1").Count(&sickCount)
	summary.TotalSickDays = int(sickCount)

	// Calculate total debt for the month
	var totalDebt int
	buildQuery(&domain.CashAdvance{}, "tanggal").Where("status = 1").Select("COALESCE(SUM(jumlah), 0)").Scan(&totalDebt)
	summary.TotalDebts = totalDebt

	return &summary, nil
}

func (r *reportRepository) GetDashboardStats(ctx context.Context, employeeID *int) (map[string]interface{}, error) {
	db := r.db.WithContext(ctx)
	stats := make(map[string]interface{})

	todayDate := time.Now().UTC().Format("2006-01-02")

	// Total Karyawan Aktif
	var totalEmployees int64
	db.Model(&domain.Employee{}).Where("id_admin_type != 1 AND status_admin = 1").Count(&totalEmployees)
	stats["total_employees"] = totalEmployees

	// Presence Hari Ini
	var presentCount int64
	qPresent := db.Model(&domain.Attendance{}).Where("tanggal = ?", todayDate)
	if employeeID != nil {
		qPresent = qPresent.Where("id_admin = ?", *employeeID)
	}
	qPresent.Count(&presentCount)

	if employeeID == nil {
		if totalEmployees > 0 {
			stats["presence_percentage"] = float64(presentCount) / float64(totalEmployees) * 100
		} else {
			stats["presence_percentage"] = 0
		}
	} else {
		stats["presence_percentage"] = presentCount // Simply 1 or 0 for user
	}

	// Sedang Cuti (Hari Ini)
	var leaveCount int64

	qLeaves := db.Model(&domain.Leave{}).
		Where("tanggal = ?", todayDate).
		Where("status = 1 AND grouping = 0")

	if employeeID != nil {
		qLeaves = qLeaves.Where("id_admin = ?", *employeeID)
	}

	qLeaves.Count(&leaveCount)
	stats["on_leave"] = leaveCount

	return stats, nil
}

func (r *reportRepository) GetPendingApprovals(ctx context.Context) (*domain.PendingApprovalsResponse, error) {
	db := r.db.WithContext(ctx)
	var resp domain.PendingApprovalsResponse
	var count int64

	db.Model(&domain.Leave{}).Where("status = 0").Count(&count)
	resp.PendingLeaves = int(count)

	db.Model(&domain.Overtime{}).Where("status = 0").Count(&count)
	resp.PendingOvertimes = int(count)

	db.Model(&domain.CashAdvance{}).Where("status = 0").Count(&count)
	resp.PendingCashAdvances = int(count)

	return &resp, nil
}

func (r *reportRepository) GetRecentActivities(ctx context.Context, employeeID *int, page, pageSize int) ([]domain.DashboardActivity, int64, error) {
	var activities []domain.DashboardActivity
	sevenDaysAgo := time.Now().UTC().AddDate(0, 0, -7)

	// We'll fetch all matching activities within the 7-day window to provide accurate total count
	// and consistent pagination. For performance, we'll cap the fetch to avoid massive memory usage
	// although 7 days of HR data should be manageable.
	maxRecent := 500

	// 1. Fetch Attendances
	var attendances []domain.Attendance
	qAttendances := r.db.WithContext(ctx).
		Model(&domain.Attendance{}).
		Select("t_presensi.*, t_admin.nama_admin as employee_name").
		Joins("JOIN t_admin ON t_admin.id_admin = t_presensi.id_admin").
		Where("t_presensi.tanggal >= ?", sevenDaysAgo)

	if employeeID != nil {
		qAttendances = qAttendances.Where("t_presensi.id_admin = ?", *employeeID)
	}

	qAttendances.Order("t_presensi.tanggal DESC, t_presensi.jam_masuk DESC").Limit(maxRecent).Find(&attendances)

	for _, a := range attendances {
		// TimeIn/TimeOut might be "HH:mm:ss" or "HH:mm"
		parseTime := func(ts string) time.Time {
			t, err := time.Parse("15:04:05", ts)
			if err != nil {
				t, _ = time.Parse("15:04", ts)
			}
			return t
		}

		// Clock In Activity
		tIn := parseTime(a.TimeIn)
		combinedIn := time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), tIn.Hour(), tIn.Minute(), tIn.Second(), 0, time.UTC)

		activities = append(activities, domain.DashboardActivity{
			ID:           a.ID,
			EmployeeName: a.EmployeeName,
			Type:         "attendance",
			Action:       "Absen Masuk",
			Date:         combinedIn,
			Status:       "Hadir",
		})

		// Clock Out Activity (if exists)
		if a.TimeOut != nil && *a.TimeOut != "" {
			tOut := parseTime(*a.TimeOut)
			combinedOut := time.Date(a.Date.Year(), a.Date.Month(), a.Date.Day(), tOut.Hour(), tOut.Minute(), tOut.Second(), 0, time.UTC)

			activities = append(activities, domain.DashboardActivity{
				ID:           a.ID,
				EmployeeName: a.EmployeeName,
				Type:         "attendance",
				Action:       "Absen Pulang",
				Date:         combinedOut,
				Status:       "Hadir",
			})
		}
	}

	// 2. Fetch Leaves
	var leaves []domain.Leave
	qLeaves := r.db.WithContext(ctx).
		Model(&domain.Leave{}).
		Select("t_cuti.*, t_admin.nama_admin as employee_name").
		Joins("JOIN t_admin ON t_admin.id_admin = t_cuti.id_admin").
		Where("t_cuti.tanggal >= ?", sevenDaysAgo)

	if employeeID != nil {
		qLeaves = qLeaves.Where("t_cuti.id_admin = ?", *employeeID)
	}

	qLeaves.Order("t_cuti.tanggal DESC").Limit(maxRecent).Find(&leaves)

	for _, l := range leaves {
		typeStr := "Cuti"
		if l.Type == 1 {
			typeStr = "Sakit"
		}
		statusStr := "Pending"
		if l.Status == 1 {
			statusStr = "Disetujui"
		} else if l.Status == 2 {
			statusStr = "Ditolak"
		}

		activities = append(activities, domain.DashboardActivity{
			ID:           l.ID,
			EmployeeName: l.EmployeeName,
			Type:         "leave",
			Action:       "Pengajuan " + typeStr,
			Date:         l.LeaveDate.UTC(),
			Status:       statusStr,
		})
	}

	// 3. Fetch Overtimes
	var overtimes []domain.Overtime
	qOvertimes := r.db.WithContext(ctx).
		Model(&domain.Overtime{}).
		Select("t_lembur.*, t_admin.nama_admin as employee_name").
		Joins("JOIN t_admin ON t_admin.id_admin = t_lembur.id_admin").
		Where("t_lembur.tanggal >= ?", sevenDaysAgo)

	if employeeID != nil {
		qOvertimes = qOvertimes.Where("t_lembur.id_admin = ?", *employeeID)
	}

	qOvertimes.Order("t_lembur.tanggal DESC").Limit(maxRecent).Find(&overtimes)

	for _, o := range overtimes {
		statusStr := "Pending"
		if o.Status == 1 {
			statusStr = "Disetujui"
		} else if o.Status == 2 {
			statusStr = "Ditolak"
		}
		// We'll use created_at for sorting overtimes to get exact chronological order
		activities = append(activities, domain.DashboardActivity{
			ID:           o.ID,
			EmployeeName: o.EmployeeName,
			Type:         "overtime",
			Action:       "Pengajuan Lembur",
			Date:         o.Date.UTC(), 
			Status:       statusStr,
		})
	}

	// 4. Fetch Cash Advances
	var cashAdvances []domain.CashAdvance
	qCashAdvances := r.db.WithContext(ctx).
		Model(&domain.CashAdvance{}).
		Select("t_kasbon.*, t_admin.nama_admin as employee_name").
		Joins("JOIN t_admin ON t_admin.id_admin = t_kasbon.id_admin").
		Where("t_kasbon.tanggal >= ?", sevenDaysAgo)

	if employeeID != nil {
		qCashAdvances = qCashAdvances.Where("t_kasbon.id_admin = ?", *employeeID)
	}

	qCashAdvances.Order("t_kasbon.tanggal DESC").Limit(maxRecent).Find(&cashAdvances)

	for _, ca := range cashAdvances {
		statusStr := "Pending"
		if ca.Status == 1 {
			statusStr = "Disetujui"
		} else if ca.Status == 2 {
			statusStr = "Ditolak"
		}
		activities = append(activities, domain.DashboardActivity{
			ID:           ca.ID,
			EmployeeName: ca.EmployeeName,
			Type:         "cash_advance",
			Action:       "Pengajuan Kasbon",
			Date:         ca.CreatedAt.UTC(),
			Status:       statusStr,
		})
	}

	// Sort merged activities by Date DESC
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Date.After(activities[j].Date)
	})

	totalCount := int64(len(activities))

	// Manual Pagination
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

	// Helper to create the base filtered query
	getQuery := func() *gorm.DB {
		base := r.db.WithContext(ctx).Model(&domain.Transaction{})

		switch filterType {
		case "monthly":
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			return base.Where("tgl_transaksi >= ? AND tgl_transaksi < ?", startDate, endDate)
		case "yearly":
			startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(1, 0, 0)
			return base.Where("tgl_transaksi >= ? AND tgl_transaksi < ?", startDate, endDate)
		default: // last-7-days
			sevenDaysAgo := now.AddDate(0, 0, -7)
			return base.Where("tgl_transaksi >= ?", sevenDaysAgo)
		}
	}

	// 1. Calculate General Stats
	var totalRevenue int64
	var totalOrders int64
	var pendingOrders int64
	var completedOrders int64
	var totalShipping int64
	var totalDiscount int64

	getQuery().Select("COALESCE(SUM(total), 0)").Scan(&totalRevenue)
	getQuery().Count(&totalOrders)
	getQuery().Where("transaksi_status = 1").Count(&pendingOrders)
	getQuery().Where("transaksi_status >= 104").Count(&completedOrders)
	getQuery().Select("COALESCE(SUM(ongkir), 0)").Scan(&totalShipping)
	getQuery().Select("COALESCE(SUM(diskon), 0)").Scan(&totalDiscount)

	stats.TotalRevenue = int(totalRevenue)
	stats.TotalOrders = int(totalOrders)
	stats.PendingOrders = int(pendingOrders)
	stats.CompletedOrders = int(completedOrders)
	stats.TotalShippingCost = int(totalShipping)
	stats.TotalDiscount = int(totalDiscount)

	// 2. Calculate Graph Data
	var graphPoints []domain.GraphDataPoint

	if filterType == "yearly" {
		for m := 1; m <= 12; m++ {
			startDate := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			var monthRev int64
			r.db.WithContext(ctx).Table("t_transaksi").
				Select("COALESCE(SUM(total), 0)").
				Where("tgl_transaksi >= ? AND tgl_transaksi < ?", startDate, endDate).
				Scan(&monthRev)

			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: startDate.Format("Jan"),
				Value: int(monthRev),
			})
		}
	} else if filterType == "monthly" {
		daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for d := 1; d <= daysInMonth; d++ {
			dayStart := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
			dayEnd := dayStart.AddDate(0, 0, 1)
			var dayRev int64
			r.db.WithContext(ctx).Table("t_transaksi").
				Select("COALESCE(SUM(total), 0)").
				Where("tgl_transaksi >= ? AND tgl_transaksi < ?", dayStart, dayEnd).
				Scan(&dayRev)

			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: dayStart.Format("02 Jan"),
				Value: int(dayRev),
			})
		}
	} else {
		for i := 6; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
			dayEnd := dayStart.AddDate(0, 0, 1)
			var dayRev int64
			r.db.WithContext(ctx).Table("t_transaksi").
				Select("COALESCE(SUM(total), 0)").
				Where("tgl_transaksi >= ? AND tgl_transaksi < ?", dayStart, dayEnd).
				Scan(&dayRev)

			graphPoints = append(graphPoints, domain.GraphDataPoint{
				Label: dayStart.Format("02 Jan"),
				Value: int(dayRev),
			})
		}
	}
	stats.RevenueGraph = graphPoints

	return &stats, nil
}

