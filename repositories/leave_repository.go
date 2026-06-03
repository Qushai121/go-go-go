package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hrms_go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaveRepository interface {
	AddCuti(data []models.LeaveHistory) error
	FindCuti(filter models.LeaveHistoryFilter) ([]models.LeaveHistory, error)
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

func (r *leaveRepository) FindCuti(filter models.LeaveHistoryFilter) ([]models.LeaveHistory, error) {
	var data []models.LeaveHistory

	query := r.db.Model(&models.LeaveHistory{})

	trim := func(value string) string {
		return strings.TrimSpace(value)
	}

	parseDate := func(value string, fieldName string) error {
		value = trim(value)
		if value == "" {
			return nil
		}

		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("format %s harus YYYY-MM-DD", fieldName)
		}

		return nil
	}

	parseInt := func(value string, fieldName string) (int, error) {
		value = trim(value)
		if value == "" {
			return 0, nil
		}

		result, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("%s harus angka", fieldName)
		}

		return result, nil
	}

	filter.LeaveHistoryId = trim(filter.LeaveHistoryId)
	filter.EmployeeNik = trim(filter.EmployeeNik)
	filter.LeaveType = trim(filter.LeaveType)
	filter.StartDate = trim(filter.StartDate)
	filter.EndDate = trim(filter.EndDate)
	filter.LeaveStart = trim(filter.LeaveStart)
	filter.LeaveEnd = trim(filter.LeaveEnd)
	filter.TotalDays = trim(filter.TotalDays)
	filter.Remarks = trim(filter.Remarks)
	filter.Location = trim(filter.Location)
	filter.LeaveYear = trim(filter.LeaveYear)
	filter.Status = trim(filter.Status)
	filter.CurrentStep = trim(filter.CurrentStep)
	filter.ApprovalHeaderId = trim(filter.ApprovalHeaderId)
	filter.ObjectCode = trim(filter.ObjectCode)
	filter.CreatedAt = trim(filter.CreatedAt)
	filter.UpdatedAt = trim(filter.UpdatedAt)
	filter.CreatedBy = trim(filter.CreatedBy)
	filter.UpdatedBy = trim(filter.UpdatedBy)
	filter.CreatedAtStart = trim(filter.CreatedAtStart)
	filter.CreatedAtEnd = trim(filter.CreatedAtEnd)
	filter.UpdatedAtStart = trim(filter.UpdatedAtStart)
	filter.UpdatedAtEnd = trim(filter.UpdatedAtEnd)

	if filter.LeaveHistoryId != "" {
		leaveHistoryId, err := uuid.Parse(filter.LeaveHistoryId)
		if err != nil {
			return nil, errors.New("leave_history_id tidak valid")
		}

		query = query.Where("leave_history_id = ?", leaveHistoryId)
	}

	if filter.EmployeeNik != "" {
		query = query.Where("employee_nik = ?", filter.EmployeeNik)
	}

	if filter.LeaveType != "" {
		query = query.Where("leave_type = ?", filter.LeaveType)
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		if err := parseDate(filter.StartDate, "start_date"); err != nil {
			return nil, err
		}

		if err := parseDate(filter.EndDate, "end_date"); err != nil {
			return nil, err
		}

		start, _ := time.Parse("2006-01-02", filter.StartDate)
		end, _ := time.Parse("2006-01-02", filter.EndDate)

		if end.Before(start) {
			return nil, errors.New("end_date tidak boleh lebih kecil dari start_date")
		}

		query = query.Where(`
			CAST(leave_start AS date) <= CAST(? AS date)
			AND CAST(leave_end AS date) >= CAST(? AS date)
		`, filter.EndDate, filter.StartDate)

	} else if filter.StartDate != "" {
		if err := parseDate(filter.StartDate, "start_date"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(leave_end AS date) >= CAST(? AS date)", filter.StartDate)

	} else if filter.EndDate != "" {
		if err := parseDate(filter.EndDate, "end_date"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(leave_start AS date) <= CAST(? AS date)", filter.EndDate)
	}

	if filter.LeaveStart != "" {
		if err := parseDate(filter.LeaveStart, "leave_start"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(leave_start AS date) = CAST(? AS date)", filter.LeaveStart)
	}

	if filter.LeaveEnd != "" {
		if err := parseDate(filter.LeaveEnd, "leave_end"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(leave_end AS date) = CAST(? AS date)", filter.LeaveEnd)
	}

	if filter.TotalDays != "" {
		totalDays, err := parseInt(filter.TotalDays, "total_days")
		if err != nil {
			return nil, err
		}

		query = query.Where("total_days = ?", totalDays)
	}

	if filter.Remarks != "" {
		query = query.Where("remarks ILIKE ?", "%"+filter.Remarks+"%")
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}

	if filter.LeaveYear != "" {
		leaveYear, err := parseInt(filter.LeaveYear, "leave_year")
		if err != nil {
			return nil, err
		}

		query = query.Where("leave_year = ?", leaveYear)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.CurrentStep != "" {
		currentStep, err := parseInt(filter.CurrentStep, "current_step")
		if err != nil {
			return nil, err
		}

		query = query.Where("current_step = ?", currentStep)
	}

	if filter.ApprovalHeaderId != "" {
		approvalHeaderId, err := uuid.Parse(filter.ApprovalHeaderId)
		if err != nil {
			return nil, errors.New("approvalheader_id tidak valid")
		}

		query = query.Where("approvalheader_id = ?", approvalHeaderId)
	}

	if filter.ObjectCode != "" {
		query = query.Where("object_code = ?", filter.ObjectCode)
	}

	if filter.CreatedAt != "" {
		if err := parseDate(filter.CreatedAt, "created_at"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(created_at AS date) = CAST(? AS date)", filter.CreatedAt)
	}

	if filter.UpdatedAt != "" {
		if err := parseDate(filter.UpdatedAt, "updated_at"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(updated_at AS date) = CAST(? AS date)", filter.UpdatedAt)
	}

	if filter.CreatedBy != "" {
		query = query.Where("created_by = ?", filter.CreatedBy)
	}

	if filter.UpdatedBy != "" {
		query = query.Where("updated_by = ?", filter.UpdatedBy)
	}

	if filter.CreatedAtStart != "" {
		if err := parseDate(filter.CreatedAtStart, "created_at_start"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(created_at AS date) >= CAST(? AS date)", filter.CreatedAtStart)
	}

	if filter.CreatedAtEnd != "" {
		if err := parseDate(filter.CreatedAtEnd, "created_at_end"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(created_at AS date) <= CAST(? AS date)", filter.CreatedAtEnd)
	}

	if filter.UpdatedAtStart != "" {
		if err := parseDate(filter.UpdatedAtStart, "updated_at_start"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(updated_at AS date) >= CAST(? AS date)", filter.UpdatedAtStart)
	}

	if filter.UpdatedAtEnd != "" {
		if err := parseDate(filter.UpdatedAtEnd, "updated_at_end"); err != nil {
			return nil, err
		}

		query = query.Where("CAST(updated_at AS date) <= CAST(? AS date)", filter.UpdatedAtEnd)
	}

	if err := query.
		Order("leave_start DESC, created_at DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("Failed to get cuti. %w", err)
	}

	return data, nil
}
