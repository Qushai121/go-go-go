package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"hrms_go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaveRepository interface {
	AddCuti(data []models.LeaveHistory) error
}

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{
		db: db,
	}
}

func (r *leaveRepository) AddCuti(data []models.LeaveHistory) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("Failed to add cuti. %w", tx.Error)
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		}
	}()

	if data == nil || len(data) == 0 {
		tx.Rollback()
		return errors.New("Failed to add cuti. Harus ada yang dikirim!")
	}

	for _, item := range data {
		employeeNik := strings.TrimSpace(item.EmployeeNik)
		leaveType := strings.ToUpper(strings.TrimSpace(item.LeaveType))

		if employeeNik == "" {
			tx.Rollback()
			return errors.New("Failed to add cuti. NIK karyawan wajib diisi.")
		}

		if leaveType == "" {
			tx.Rollback()
			return errors.New("Failed to add cuti. Tipe cuti wajib diisi.")
		}

		if item.LeaveStart.IsZero() {
			tx.Rollback()
			return errors.New("Failed to add cuti. Tanggal mulai cuti wajib diisi.")
		}

		if item.LeaveEnd.IsZero() {
			tx.Rollback()
			return errors.New("Failed to add cuti. Tanggal akhir cuti wajib diisi.")
		}

		if dateOnly(item.LeaveEnd).Before(dateOnly(item.LeaveStart)) {
			tx.Rollback()
			return errors.New("Failed to add cuti. Tanggal akhir cuti tidak boleh lebih kecil dari tanggal mulai.")
		}

		totalDays := item.TotalDays

		if leaveType == "CUTI" {
			calcTotalDays, err := calculateLeaveTotalDays(tx, employeeNik, item.LeaveStart, item.LeaveEnd)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to add cuti. %w", err)
			}

			totalDays = calcTotalDays
		}

		if leaveType != "CUTI" && totalDays <= 0 {
			totalDays = calculateCalendarDays(item.LeaveStart, item.LeaveEnd)
		}

		leaveYear := item.LeaveYear
		if leaveYear <= 0 {
			leaveYear = item.LeaveStart.Year()
		}

		currentStep := item.CurrentStep
		if currentStep <= 0 {
			currentStep = 1
		}

		objectCode := strings.TrimSpace(item.ObjectCode)
		if objectCode == "" {
			objectCode = "LEAVE_HISTORY"
		}

		createdAt := item.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}

		createdBy := strings.TrimSpace(item.CreatedBy)
		if createdBy == "" {
			createdBy = "System"
		}

		leaveHistory := models.LeaveHistory{
			LeaveHistoryId:   uuid.New(),
			EmployeeNik:      employeeNik,
			LeaveType:        leaveType,
			LeaveStart:       item.LeaveStart,
			LeaveEnd:         item.LeaveEnd,
			TotalDays:        totalDays,
			Remarks:          item.Remarks,
			Location:         item.Location,
			LeaveYear:        leaveYear,
			Status:           "A",
			CurrentStep:      currentStep,
			ApprovalHeaderId: item.ApprovalHeaderId,
			ObjectCode:       objectCode,
			CreatedAt:        createdAt,
			CreatedBy:        createdBy,
		}

		if leaveHistory.LeaveType == "CUTI" {
			calendarDays := calculateCalendarDays(leaveHistory.LeaveStart, leaveHistory.LeaveEnd)

			if calendarDays <= 0 {
				tx.Rollback()
				return errors.New("Failed to add cuti. Total hari cuti tidak valid.")
			}

			err := tx.Exec(`
				CALL public.sp_deduct_leave(
					CAST(? AS varchar),
					CAST(? AS date),
					CAST(? AS integer),
					CAST(? AS varchar)
				)
			`,
				leaveHistory.EmployeeNik,
				leaveHistory.LeaveStart.Format("2006-01-02"),
				calendarDays,
				leaveHistory.CreatedBy,
			).Error

			if err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to add cuti. %w", err)
			}
		}

		if err := tx.Create(&leaveHistory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to add cuti. %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to add cuti. %w", err)
	}

	return nil
}

func calculateLeaveTotalDays(tx *gorm.DB, employeeNik string, leaveStart time.Time, leaveEnd time.Time) (int, error) {
	_ = employeeNik

	startDate := dateOnly(leaveStart)
	endDate := dateOnly(leaveEnd)

	if endDate.Before(startDate) {
		return 0, errors.New("tanggal akhir tidak boleh lebih kecil dari tanggal mulai")
	}

	var totalDays int

	err := tx.Raw(`
		SELECT COALESCE(COUNT(1), 0)::int
		FROM public.hrms_wkcal
		WHERE date_id BETWEEN CAST(? AS date) AND CAST(? AS date)
		  AND COALESCE(is_working_day, false) = true
	`,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	).Scan(&totalDays).Error

	if err != nil {
		return 0, err
	}

	return totalDays, nil
}

func calculateCalendarDays(leaveStart time.Time, leaveEnd time.Time) int {
	if leaveStart.IsZero() || leaveEnd.IsZero() {
		return 0
	}

	startDate := dateOnly(leaveStart)
	endDate := dateOnly(leaveEnd)

	if endDate.Before(startDate) {
		return 0
	}

	return int(endDate.Sub(startDate).Hours()/24) + 1
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}
