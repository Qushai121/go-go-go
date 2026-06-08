package repositories

import (
	"fmt"
	"strings"

	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

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

	if err := r.normalizeAndValidateSite(attendance); err != nil {
		return err
	}

	err := r.db.Model(&models.Attendance{}).
		Where("user_id = ? AND company_code = ? AND site_type = ? AND site_code = ? AND logtime = ? AND functionno = ?",
			attendance.UserId,
			attendance.CompanyCode,
			attendance.SiteType,
			attendance.SiteCode,
			attendance.LogTime,
			attendance.FunctionNo).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("attendance already exists for this site and logtime")
	}

	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) normalizeAndValidateSite(attendance *models.Attendance) error {
	attendance.SiteType = strings.ToUpper(strings.TrimSpace(attendance.SiteType))
	attendance.SiteCode = strings.TrimSpace(attendance.SiteCode)
	attendance.CompanyCode = strings.TrimSpace(attendance.CompanyCode)

	if attendance.SiteCode == "" {
		return fmt.Errorf("site_code is required")
	}

	query := r.db.Model(&models.Site{}).Where("site_code = ?", attendance.SiteCode)
	if attendance.CompanyCode != "" {
		query = query.Where("company_code = ?", attendance.CompanyCode)
	}
	if attendance.SiteType != "" {
		query = query.Where("site_type = ?", attendance.SiteType)
	}

	var site models.Site
	if err := query.Limit(1).Take(&site).Error; err != nil {
		return fmt.Errorf("site not found: %w", err)
	}

	if attendance.CompanyCode == "" {
		attendance.CompanyCode = site.CompanyCode
	}
	if attendance.SiteType == "" {
		attendance.SiteType = site.SiteType
	}
	if attendance.MaxRadius == nil && site.MaxRadius > 0 {
		attendance.MaxRadius = &site.MaxRadius
	}

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
		args := make([]interface{}, 22)
		for i := range args {
			args[i] = search
		}

		modelDb = modelDb.Where(`
			attendance_id::text ILIKE ? OR
			user_id::text ILIKE ? OR
			company_code ILIKE ? OR
			site_type ILIKE ? OR
			site_code ILIKE ? OR
			logtime::text ILIKE ? OR
			functionno::text ILIKE ? OR
			activity_type ILIKE ? OR
			latitude ILIKE ? OR
			longitude ILIKE ? OR
			presentase_kemiripan ILIKE ? OR
			imagepath ILIKE ? OR
			is_offline ILIKE ? OR
			distance ILIKE ? OR
			platforms ILIKE ? OR
			max_radius::text ILIKE ? OR
			expand_radius::text ILIKE ? OR
			object_code ILIKE ? OR
			created_at::text ILIKE ? OR
			updated_at::text ILIKE ? OR
			created_by ILIKE ? OR
			updated_by ILIKE ?
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
			COALESCE(a.site_type, '') AS site_type,
			COALESCE(a.site_code, '') AS site_code,
			COALESCE(s.site_name, '') AS site_name,
			COALESCE(s.site_phone, '') AS site_phone,
			COALESCE(s.site_address, '') AS site_address,
			COALESCE(s.site_latitude, '') AS site_latitude,
			COALESCE(s.site_longitude, '') AS site_longitude,
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
			LEFT JOIN hrms_site s
				ON s.company_code = a.company_code
				AND s.site_type = a.site_type
				AND s.site_code = a.site_code
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
			a.site_type ILIKE ? OR
			a.site_code ILIKE ? OR
			s.site_name ILIKE ? OR
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
		`, search, search, search, search, search, search, search, search, search, search, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := attendanceMeDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func attendanceMeDynamicSearchFields() map[string]dto.DynamicSearchDto {
	return map[string]dto.DynamicSearchDto{
		"attendance_id": {Field: "a.attendance_id", Query: " = ?"},
		"user_id":       {Field: "a.user_id", Query: " = ?"},
		"company_code":  {Field: "a.company_code", Query: " ILIKE ?"},
		"site_type":     {Field: "a.site_type", Query: " ILIKE ?"},
		"site_code":     {Field: "a.site_code", Query: " ILIKE ?"},
		"site_name":     {Field: "s.site_name", Query: " ILIKE ?"},
		"logtime":       {Field: "a.logtime", Query: " >= ?"},
		"functionno":    {Field: "a.functionno", Query: " = ?"},
		"activity_type": {Field: "a.activity_type", Query: " ILIKE ?"},
		"action_type":   {Field: "a.action_type", Query: " ILIKE ?"},
	}
}

func attendanceDynamicSearchFields() map[string]dto.DynamicSearchDto {
	return map[string]dto.DynamicSearchDto{
		"attendance_id":        {Field: "attendance_id", Query: " = ?"},
		"user_id":              {Field: "user_id", Query: " = ?"},
		"company_code":         {Field: "company_code", Query: " ILIKE ?"},
		"site_type":            {Field: "site_type", Query: " ILIKE ?"},
		"site_code":            {Field: "site_code", Query: " ILIKE ?"},
		"logtime":              {Field: "logtime", Query: " >= ?"},
		"functionno":           {Field: "functionno", Query: " = ?"},
		"activity_type":        {Field: "activity_type", Query: " ILIKE ?"},
		"latitude":             {Field: "latitude", Query: " ILIKE ?"},
		"longitude":            {Field: "longitude", Query: " ILIKE ?"},
		"presentase_kemiripan": {Field: "presentase_kemiripan", Query: " ILIKE ?"},
		"imagepath":            {Field: "imagepath", Query: " ILIKE ?"},
		"is_offline":           {Field: "is_offline", Query: " ILIKE ?"},
		"distance":             {Field: "distance", Query: " ILIKE ?"},
		"platforms":            {Field: "platforms", Query: " ILIKE ?"},
		"max_radius":           {Field: "max_radius", Query: " = ?"},
		"expand_radius":        {Field: "expand_radius", Query: " = ?"},
		"object_code":          {Field: "object_code", Query: " ILIKE ?"},
		"created_by":           {Field: "created_by", Query: " ILIKE ?"},
		"updated_by":           {Field: "updated_by", Query: " ILIKE ?"},
	}
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db}
}
