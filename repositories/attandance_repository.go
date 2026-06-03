package repositories

import (
	"fmt"
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"
	"strings"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error)
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]response.AttendanceMeResponseDto], error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	var count int64

	if err := r.fillBranchCodeFromOffice(attendance); err != nil {
		return err
	}

	err := r.db.Model(&models.Attendance{}).
		Where("user_id = ? AND company_code = ? AND office_code = ? AND logtime = ? AND functionno = ?",
			attendance.UserId,
			attendance.CompanyCode,
			attendance.OfficeCode,
			attendance.LogTime,
			attendance.FunctionNo).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("attendance already exists for this logtime")
	}

	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) fillBranchCodeFromOffice(attendance *models.Attendance) error {
	if strings.TrimSpace(attendance.BranchCode) != "" {
		return nil
	}

	if strings.TrimSpace(attendance.OfficeCode) == "" {
		return fmt.Errorf("office_code is required")
	}

	type officeBranchResult struct {
		BranchCode string `gorm:"column:branch_code"`
	}

	var result officeBranchResult

	query := r.db.
		Table("hrms_office").
		Select("branch_code").
		Where("office_code = ?", attendance.OfficeCode)

	if strings.TrimSpace(attendance.CompanyCode) != "" {
		query = query.Where("company_code = ?", attendance.CompanyCode)
	}

	err := query.Limit(1).Take(&result).Error
	if err != nil {
		return fmt.Errorf("failed to get branch_code from office: %w", err)
	}

	if strings.TrimSpace(result.BranchCode) == "" {
		return fmt.Errorf("branch_code not found for office_code %s", attendance.OfficeCode)
	}

	attendance.BranchCode = result.BranchCode
	return nil
}

func (r *attendanceRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error) {
	var data []models.Attendance
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Attendance{}).Preload("User")

	dataAkhir := response.PaginateResponseDto[[]models.Attendance]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "attendance_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		args := make([]interface{}, 23)
		for i := range args {
			args[i] = search
		}

		modelDb = modelDb.Where(`
			attendance_id::text LIKE ? OR
			user_id::text LIKE ? OR
			company_code LIKE ? OR
			office_code LIKE ? OR
			branch_code LIKE ? OR
			customer_code LIKE ? OR
			logtime::text LIKE ? OR
			functionno::text LIKE ? OR
			activity_type LIKE ? OR
			latitude LIKE ? OR
			longitude LIKE ? OR
			presentase_kemiripan LIKE ? OR
			imagepath LIKE ? OR
			is_offline LIKE ? OR
			distance LIKE ? OR
			platforms LIKE ? OR
			max_radius::text LIKE ? OR
			expand_radius::text LIKE ? OR
			object_code LIKE ? OR
			created_at::text LIKE ? OR
			updated_at::text LIKE ? OR
			created_by LIKE ? OR
			updated_by LIKE ?
		`, args...)
	}

	allowedDynamicList := attendanceDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func (r *attendanceRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]response.AttendanceMeResponseDto], error) {
	var data []response.AttendanceMeResponseDto
	var totalRecord int64
	var totalPage int

	modelDb := r.db.
		Table("hrms_attendance AS a").
		Select(`
			a.attendance_id::text AS attendance_id,
			a.user_id::text AS user_id,
			COALESCE(a.company_code, '') AS company_code,
			COALESCE(a.branch_code, '') AS branch_code,
			COALESCE(b.branch_name, '') AS branch_name,
			COALESCE(a.office_code, '') AS office_code,
			COALESCE(o.office_name, '') AS office_name,
			COALESCE(a.customer_code, '') AS customer_code,
			COALESCE(c.customer_name, '') AS customer_name,
			a.logtime AS logtime,
			a.functionno AS functionno,
			COALESCE(a.activity_type, '') AS activity_type,
			COALESCE(a.latitude, '') AS latitude,
			COALESCE(a.longitude, '') AS longitude,
			COALESCE(a.presentase_kemiripan, '') AS presentase_kemiripan,
			COALESCE(a.imagepath, '') AS imagepath,
			COALESCE(a.is_offline, '') AS is_offline,
			COALESCE(a.distance, '') AS distance,
			COALESCE(a.platforms, '') AS platforms,
			COALESCE(a.max_radius::text, '') AS max_radius,
			COALESCE(a.expand_radius::text, '') AS expand_radius,
			COALESCE(a.object_code, '') AS object_code,
			a.created_at AS created_at,
			a.updated_at AS updated_at,
			COALESCE(a.created_by, '') AS created_by,
			COALESCE(a.updated_by, '') AS updated_by
		`).
		Joins(`
			LEFT JOIN hrms_branch b
				ON b.company_code = a.company_code
				AND b.branch_code = a.branch_code
		`).
		Joins(`
			LEFT JOIN hrms_office o
				ON o.company_code = a.company_code
				AND o.branch_code = a.branch_code
				AND o.office_code = a.office_code
		`).
		Joins(`
			LEFT JOIN hrms_customer c
				ON c.customer_code = a.customer_code
		`).
		Where("a.user_id = ?", userId)

	dataAkhir := response.PaginateResponseDto[[]response.AttendanceMeResponseDto]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "a.logtime"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			a.attendance_id::text ILIKE ? OR
			a.user_id::text ILIKE ? OR
			a.company_code ILIKE ? OR
			a.branch_code ILIKE ? OR
			b.branch_name ILIKE ? OR
			a.office_code ILIKE ? OR
			o.office_name ILIKE ? OR
			a.customer_code ILIKE ? OR
			c.customer_name ILIKE ? OR
			a.logtime::text ILIKE ? OR
			a.functionno::text ILIKE ? OR
			a.activity_type ILIKE ? OR
			a.latitude ILIKE ? OR
			a.longitude ILIKE ? OR
			a.presentase_kemiripan ILIKE ? OR
			a.imagepath ILIKE ? OR
			a.is_offline ILIKE ? OR
			a.distance ILIKE ? OR
			a.platforms ILIKE ? OR
			a.object_code ILIKE ? OR
			a.created_by ILIKE ? OR
			a.updated_by ILIKE ?
		`,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search, search,
		)
	}

	allowedDynamicList := attendanceMeDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func attendanceMeDynamicSearchFields() map[string]dto.DynamicSearchDto {
	return map[string]dto.DynamicSearchDto{
		"attendance_id": {
			Field: "a.attendance_id",
			Query: " = ?",
		},
		"user_id": {
			Field: "a.user_id",
			Query: " = ?",
		},
		"company_code": {
			Field: "a.company_code",
			Query: " ILIKE ?",
		},
		"branch_code": {
			Field: "a.branch_code",
			Query: " ILIKE ?",
		},
		"branch_name": {
			Field: "b.branch_name",
			Query: " ILIKE ?",
		},
		"office_code": {
			Field: "a.office_code",
			Query: " ILIKE ?",
		},
		"office_name": {
			Field: "o.office_name",
			Query: " ILIKE ?",
		},
		"customer_code": {
			Field: "a.customer_code",
			Query: " ILIKE ?",
		},
		"customer_name": {
			Field: "c.customer_name",
			Query: " ILIKE ?",
		},
		"logtime": {
			Field: "a.logtime",
			Query: " >= ?",
		},
		"functionno": {
			Field: "a.functionno",
			Query: " = ?",
		},
		"activity_type": {
			Field: "a.activity_type",
			Query: " ILIKE ?",
		},
		"action_type": {
			Field: "a.action_type",
			Query: " ILIKE ?",
		},
	}
}

func attendanceDynamicSearchFields() map[string]dto.DynamicSearchDto {
	return map[string]dto.DynamicSearchDto{
		"attendance_id": {
			Field: "attendance_id",
			Query: " = ?",
		},
		"user_id": {
			Field: "user_id",
			Query: " = ?",
		},
		"company_code": {
			Field: "company_code",
			Query: " LIKE ?",
		},
		"office_code": {
			Field: "office_code",
			Query: " LIKE ?",
		},
		"branch_code": {
			Field: "branch_code",
			Query: " LIKE ?",
		},
		"customer_code": {
			Field: "customer_code",
			Query: " LIKE ?",
		},
		"logtime": {
			Field: "logtime",
			Query: " >= ?",
		},
		"functionno": {
			Field: "functionno",
			Query: " = ?",
		},
		"activity_type": {
			Field: "activity_type",
			Query: " LIKE ?",
		},
		"latitude": {
			Field: "latitude",
			Query: " LIKE ?",
		},
		"longitude": {
			Field: "longitude",
			Query: " LIKE ?",
		},
		"presentase_kemiripan": {
			Field: "presentase_kemiripan",
			Query: " LIKE ?",
		},
		"imagepath": {
			Field: "imagepath",
			Query: " LIKE ?",
		},
		"is_offline": {
			Field: "is_offline",
			Query: " LIKE ?",
		},
		"distance": {
			Field: "distance",
			Query: " LIKE ?",
		},
		"platforms": {
			Field: "platforms",
			Query: " LIKE ?",
		},
		"max_radius": {
			Field: "max_radius",
			Query: " = ?",
		},
		"expand_radius": {
			Field: "expand_radius",
			Query: " = ?",
		},
		"object_code": {
			Field: "object_code",
			Query: " LIKE ?",
		},
		"created_by": {
			Field: "created_by",
			Query: " LIKE ?",
		},
		"updated_by": {
			Field: "updated_by",
			Query: " LIKE ?",
		},
	}
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db}
}
