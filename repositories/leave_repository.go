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
	AddCuti(data models.LeaveHistory) error
	FindCuti(filter models.LeaveHistoryFilter) ([]models.LeaveHistory, error)
	GetBalance(employeeNik string, leaveTypeId *uuid.UUID) (*models.LeaveBalanceResponse, error)
	CreateTransaction(request models.CreateLeaveTransactionRequest) (*models.CreateLeaveTransactionResponse, error)
}

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{
		db: db,
	}
}

func (r *leaveRepository) AddCuti(item models.LeaveHistory) error {
	// Support logic cuti dynamic tanpa memutus endpoint lama /api/leave/cuti.
	// Jika request sudah membawa leave_type_id, prosesnya dialihkan ke CreateTransaction
	// sehingga range tanggal tetap dipecah menjadi transaksi per tanggal.
	if item.LeaveTypeId != nil && *item.LeaveTypeId != uuid.Nil {
		createdBy := strings.TrimSpace(item.CreatedBy)
		if createdBy == "" {
			createdBy = item.EmployeeNik
		}

		_, err := r.CreateTransaction(models.CreateLeaveTransactionRequest{
			EmployeeNik:      item.EmployeeNik,
			LeaveTypeId:      item.LeaveTypeId.String(),
			LeaveStart:       item.LeaveStart.Format("2006-01-02"),
			LeaveEnd:         item.LeaveEnd.Format("2006-01-02"),
			Remarks:          stringValue(item.Remarks),
			Location:         stringValue(item.Location),
			CurrentStep:      item.CurrentStep,
			ApprovalHeaderId: uuidValue(item.ApprovalHeaderId),
			CreatedBy:        createdBy,
		})

		return err
	}

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
		LeaveDate:        dateOnly(item.LeaveStart),
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
	filter.LeaveTypeId = trim(filter.LeaveTypeId)
	filter.LeaveType = trim(filter.LeaveType)
	filter.StartDate = trim(filter.StartDate)
	filter.EndDate = trim(filter.EndDate)
	filter.LeaveDate = trim(filter.LeaveDate)
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

	if filter.LeaveTypeId != "" {
		leaveTypeId, err := uuid.Parse(filter.LeaveTypeId)
		if err != nil {
			return nil, errors.New("leave_type_id tidak valid")
		}

		query = query.Where("leave_type_id = ?", leaveTypeId)
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

	if filter.LeaveDate != "" {
		if err := parseDate(filter.LeaveDate, "leave_date"); err != nil {
			return nil, err
		}

		query = query.Where("COALESCE(leave_date::date, leave_start::date) = CAST(? AS date)", filter.LeaveDate)
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

func (r *leaveRepository) GetBalance(employeeNik string, leaveTypeId *uuid.UUID) (*models.LeaveBalanceResponse, error) {
	employeeNik = strings.TrimSpace(employeeNik)
	if employeeNik == "" {
		return nil, errors.New("NIK karyawan wajib diisi")
	}

	query := r.db.Table("hrms_leave_balance b").
		Select(`
			b.leave_balance_id,
			b.employee_nik,
			b.leave_type_id,
			lt.leave_type_code,
			lt.leave_type_name,
			b.leave_period_start,
			b.leave_period_end,
			COALESCE(b.total_leave, 0)::int AS total_leave,
			COALESCE(b.leave_used, 0)::int AS leave_used,
			COALESCE(b.carry_forward, 0)::int AS carry_forward,
			(
				COALESCE(b.total_leave, 0)
				+ COALESCE(b.carry_forward, 0)
				- COALESCE(b.leave_used, 0)
			)::int AS leave_remaining,
			b.object_code,
			b.created_at,
			b.updated_at,
			b.created_by,
			b.updated_by
		`).
		Joins("LEFT JOIN hrms_leave_type lt ON lt.leave_type_id = b.leave_type_id").
		Where("b.employee_nik = ?", employeeNik).
		Where("CURRENT_DATE BETWEEN b.leave_period_start::date AND b.leave_period_end::date")

	if leaveTypeId != nil {
		query = query.Where("b.leave_type_id = ?", *leaveTypeId)
	}

	var balances []models.LeaveBalance
	if err := query.Order("lt.sort_order ASC, b.leave_period_start DESC").Find(&balances).Error; err != nil {
		return nil, err
	}

	result := &models.LeaveBalanceResponse{
		EmployeeNik: employeeNik,
		Balances:    balances,
	}

	for _, item := range balances {
		result.TotalRemaining += item.LeaveRemaining

		periodStart := item.LeavePeriodStart.Format("2006-01-02")
		periodEnd := item.LeavePeriodEnd.Format("2006-01-02")

		if result.PeriodStart == "" || periodStart < result.PeriodStart {
			result.PeriodStart = periodStart
		}

		if result.PeriodEnd == "" || periodEnd > result.PeriodEnd {
			result.PeriodEnd = periodEnd
		}
	}

	return result, nil
}

func (r *leaveRepository) CreateTransaction(request models.CreateLeaveTransactionRequest) (*models.CreateLeaveTransactionResponse, error) {
	request.EmployeeNik = strings.TrimSpace(request.EmployeeNik)
	request.LeaveTypeId = strings.TrimSpace(request.LeaveTypeId)
	request.CreatedBy = strings.TrimSpace(request.CreatedBy)
	request.ApprovalHeaderId = strings.TrimSpace(request.ApprovalHeaderId)

	if request.EmployeeNik == "" {
		return nil, errors.New("NIK karyawan wajib diisi")
	}

	leaveTypeId, err := uuid.Parse(request.LeaveTypeId)
	if err != nil {
		return nil, errors.New("leave_type_id tidak valid")
	}

	leaveStart, err := parseLeaveRequestDate(request.LeaveStart, "leave_start")
	if err != nil {
		return nil, err
	}

	leaveEnd, err := parseLeaveRequestDate(request.LeaveEnd, "leave_end")
	if err != nil {
		return nil, err
	}

	if leaveEnd.Before(leaveStart) {
		return nil, errors.New("leave_end tidak boleh lebih kecil dari leave_start")
	}

	createdBy := request.CreatedBy
	if createdBy == "" {
		createdBy = request.EmployeeNik
	}

	var result *models.CreateLeaveTransactionResponse

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureEmployeeExists(tx, request.EmployeeNik); err != nil {
			return err
		}

		leaveType, err := getLeaveTypeByID(tx, leaveTypeId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("tipe cuti tidak ditemukan")
			}
			return err
		}

		if !leaveType.IsActive {
			return errors.New("tipe cuti tidak aktif")
		}

		if leaveType.MinNoticeDays > 0 {
			minDate := dateOnly(jakartaNow()).AddDate(0, 0, leaveType.MinNoticeDays)
			if leaveStart.Before(minDate) {
				return fmt.Errorf("pengajuan %s minimal H-%d", leaveType.LeaveTypeName, leaveType.MinNoticeDays)
			}
		}

		leaveDates, err := getDynamicLeaveDates(tx, request.EmployeeNik, leaveStart, leaveEnd, leaveType.UseWorkingDay)
		if err != nil {
			return err
		}

		if len(leaveDates) == 0 {
			return errors.New("tidak ada tanggal cuti yang valid pada range tersebut")
		}

		if leaveType.MaxDaysPerRequest != nil && *leaveType.MaxDaysPerRequest > 0 && len(leaveDates) > *leaveType.MaxDaysPerRequest {
			return fmt.Errorf("maksimal pengajuan %s adalah %d hari", leaveType.LeaveTypeName, *leaveType.MaxDaysPerRequest)
		}

		if err := ensureNoDuplicateLeave(tx, request.EmployeeNik, leaveDates); err != nil {
			return err
		}

		deductedDays := 0
		if leaveType.DeductLeaveBalance {
			if err := deductDynamicLeaveBalance(tx, request.EmployeeNik, leaveTypeId, len(leaveDates), createdBy); err != nil {
				return err
			}
			deductedDays = len(leaveDates)
		}

		approvalHeaderId, err := parseOptionalUUID(request.ApprovalHeaderId, "approvalheader_id")
		if err != nil {
			return err
		}

		remarks := nullableString(request.Remarks)
		location := nullableString(request.Location)
		currentStep := request.CurrentStep
		if currentStep <= 0 {
			currentStep = 1
		}

		createdAt := jakartaNow()
		leaveTypeIdCopy := leaveTypeId
		transactions := make([]models.LeaveHistory, 0, len(leaveDates))
		insertedDates := make([]string, 0, len(leaveDates))

		for _, leaveDate := range leaveDates {
			leaveDate = dateOnly(leaveDate)
			transactions = append(transactions, models.LeaveHistory{
				LeaveHistoryId:   uuid.New(),
				EmployeeNik:      request.EmployeeNik,
				LeaveTypeId:      &leaveTypeIdCopy,
				LeaveType:        leaveType.LeaveTypeCode,
				LeaveDate:        leaveDate,
				LeaveStart:       leaveDate,
				LeaveEnd:         leaveDate,
				TotalDays:        1,
				Remarks:          remarks,
				Location:         location,
				LeaveYear:        leaveDate.Year(),
				Status:           "A",
				CurrentStep:      currentStep,
				ApprovalHeaderId: approvalHeaderId,
				ObjectCode:       "LEAVE_HISTORY",
				CreatedAt:        createdAt,
				CreatedBy:        createdBy,
			})

			insertedDates = append(insertedDates, leaveDate.Format("2006-01-02"))
		}

		if err := tx.Create(&transactions).Error; err != nil {
			return fmt.Errorf("gagal insert transaksi cuti: %w", err)
		}

		result = &models.CreateLeaveTransactionResponse{
			EmployeeNik:       request.EmployeeNik,
			LeaveTypeId:       leaveTypeId.String(),
			LeaveTypeCode:     leaveType.LeaveTypeCode,
			LeaveTypeName:     leaveType.LeaveTypeName,
			LeaveStart:        leaveStart.Format("2006-01-02"),
			LeaveEnd:          leaveEnd.Format("2006-01-02"),
			InsertedDays:      len(transactions),
			DeductedDays:      deductedDays,
			InsertedDates:     insertedDates,
			LeaveTransactions: transactions,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func ensureEmployeeExists(tx *gorm.DB, employeeNik string) error {
	var total int64
	if err := tx.Table("hrms_users").Where("employee_nik = ?", employeeNik).Count(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		return errors.New("karyawan tidak ditemukan")
	}

	return nil
}

func getLeaveTypeByID(tx *gorm.DB, leaveTypeId uuid.UUID) (*models.LeaveType, error) {
	var data models.LeaveType
	err := tx.Where("leave_type_id = ?", leaveTypeId).First(&data).Error
	return &data, err
}

type dynamicLeaveDateRow struct {
	LeaveDate time.Time `gorm:"column:leave_date"`
}

func getDynamicLeaveDates(tx *gorm.DB, employeeNik string, leaveStart time.Time, leaveEnd time.Time, useWorkingDay bool) ([]time.Time, error) {
	var rows []dynamicLeaveDateRow

	err := tx.Raw(`
		SELECT d.leave_day::date AS leave_date
		FROM generate_series(CAST(? AS date), CAST(? AS date), interval '1 day') AS d(leave_day)
		WHERE
		(
			CAST(? AS boolean) = false
			OR COALESCE((
				SELECT w.is_working_day
				FROM hrms_wkcal w
				WHERE w.date_id::date = d.leave_day::date
				LIMIT 1
			), false) = true
		)
		AND
		(
			CAST(? AS boolean) = false
			OR EXISTS
			(
				SELECT 1
				FROM hrms_employee_shift es
				INNER JOIN hrms_employee_shift_period esp
					ON esp.weekly_id = es.weekly_id
				WHERE es.employee_nik = ?
				  AND es.weekday_id = EXTRACT(ISODOW FROM d.leave_day)::int
				  AND d.leave_day::date BETWEEN esp.week_start_date::date AND esp.week_end_date::date
			)
		)
		ORDER BY d.leave_day::date ASC
	`,
		leaveStart.Format("2006-01-02"),
		leaveEnd.Format("2006-01-02"),
		useWorkingDay,
		useWorkingDay,
		employeeNik,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]time.Time, 0, len(rows))
	for _, row := range rows {
		result = append(result, dateOnly(row.LeaveDate))
	}

	return result, nil
}

func ensureNoDuplicateLeave(tx *gorm.DB, employeeNik string, dates []time.Time) error {
	if len(dates) == 0 {
		return nil
	}

	dateStrings := make([]string, 0, len(dates))
	for _, date := range dates {
		dateStrings = append(dateStrings, dateOnly(date).Format("2006-01-02"))
	}

	var total int64
	if err := tx.Table("hrms_leave_history").
		Where("employee_nik = ?", employeeNik).
		Where("COALESCE(status, '') NOT IN ('C', 'R')").
		Where(`
			COALESCE(leave_date::date, leave_start::date) IN ?
			OR EXISTS (
				SELECT 1
				FROM generate_series(leave_start::date, leave_end::date, interval '1 day') AS x(dt)
				WHERE x.dt::date IN ?
			)
		`, dateStrings, dateStrings).
		Count(&total).Error; err != nil {
		return err
	}

	if total > 0 {
		return errors.New("sudah ada transaksi cuti pada salah satu tanggal tersebut")
	}

	return nil
}

func deductDynamicLeaveBalance(tx *gorm.DB, employeeNik string, leaveTypeId uuid.UUID, deductDays int, updatedBy string) error {
	var balance models.LeaveBalance

	err := tx.Raw(`
		SELECT
			leave_balance_id,
			employee_nik,
			leave_type_id,
			leave_period_start,
			leave_period_end,
			COALESCE(total_leave, 0)::int AS total_leave,
			COALESCE(leave_used, 0)::int AS leave_used,
			COALESCE(carry_forward, 0)::int AS carry_forward,
			(
				COALESCE(total_leave, 0)
				+ COALESCE(carry_forward, 0)
				- COALESCE(leave_used, 0)
			)::int AS leave_remaining
		FROM hrms_leave_balance
		WHERE employee_nik = ?
		  AND leave_type_id = ?
		  AND CURRENT_DATE BETWEEN leave_period_start::date AND leave_period_end::date
		ORDER BY leave_period_start DESC
		LIMIT 1
		FOR UPDATE
	`, employeeNik, leaveTypeId).Scan(&balance).Error
	if err != nil {
		return err
	}

	if balance.LeaveBalanceId == uuid.Nil {
		return errors.New("saldo cuti tidak ditemukan untuk tipe cuti ini")
	}

	if balance.LeaveRemaining < deductDays {
		return fmt.Errorf("saldo cuti tidak mencukupi. Sisa saldo: %d, kebutuhan: %d", balance.LeaveRemaining, deductDays)
	}

	if err := tx.Model(&models.LeaveBalance{}).
		Where("leave_balance_id = ?", balance.LeaveBalanceId).
		Updates(map[string]interface{}{
			"leave_used": gorm.Expr("COALESCE(leave_used, 0) + ?", deductDays),
			"updated_by": updatedBy,
			"updated_at": jakartaNow(),
		}).Error; err != nil {
		return err
	}

	return nil
}

func parseLeaveRequestDate(value string, fieldName string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("%s wajib diisi", fieldName)
	}

	parsed, err := time.ParseInLocation("2006-01-02", value, jakartaLocation())
	if err != nil {
		return time.Time{}, fmt.Errorf("format %s harus YYYY-MM-DD", fieldName)
	}

	return dateOnly(parsed), nil
}

func parseOptionalUUID(value string, fieldName string) (*uuid.UUID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("%s tidak valid", fieldName)
	}

	return &parsed, nil
}

func nullableString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func uuidValue(value *uuid.UUID) string {
	if value == nil || *value == uuid.Nil {
		return ""
	}

	return value.String()
}

func jakartaNow() time.Time {
	return time.Now().In(jakartaLocation())
}

func jakartaLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("Asia/Jakarta", 7*60*60)
	}

	return loc
}
