package repository

import (
	"context"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
	"gorm.io/gorm"
)

type attendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository will create an object that requests the domain.AttendanceRepository interface
func NewAttendanceRepository(db *gorm.DB) domain.AttendanceRepository {
	return &attendanceRepository{db}
}

func (r *attendanceRepository) Fetch(ctx context.Context, filter map[string]interface{}) ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	query := r.db.WithContext(ctx).Model(&domain.Attendance{})

	// Sanitize and apply filters explicitly (Senior BE approach)
	if empID, ok := filter["employee_id"]; ok {
		query = query.Where("t_presensi.id_admin = ?", empID)
	}
	if date, ok := filter["presence_date"]; ok {
		query = query.Where("t_presensi.tanggal = ?", date)
	}

	err := query.Select("t_presensi.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_presensi.id_admin").
		Find(&attendances).Error

	return attendances, err
}

func (r *attendanceRepository) GetByDateAndEmployee(ctx context.Context, date time.Time, employeeID int) (*domain.Attendance, error) {
	var att domain.Attendance
	err := r.db.WithContext(ctx).
		Where("tanggal = ? AND id_admin = ?", date.Format("2006-01-02"), employeeID).
		First(&att).Error

	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *attendanceRepository) GetTodayAttendances(ctx context.Context) ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	today := time.Now().Format("2006-01-02")

	err := r.db.WithContext(ctx).Model(&domain.Attendance{}).
		Select("t_presensi.*, t_admin.nama_admin as employee_name").
		Joins("LEFT JOIN t_admin ON t_admin.id_admin = t_presensi.id_admin").
		Where("tanggal = ?", today).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *attendanceRepository) Store(ctx context.Context, attendance *domain.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *domain.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}
