package repositories

import (
	"fmt"
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error)
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	var count int64

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

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			attendance_id::text LIKE ? OR
			user_id::text LIKE ? OR
			company_code LIKE ? OR
			office_code LIKE ? OR
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
		`,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search,
		)
	}

	allowedDynamicList := attendanceDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (r *attendanceRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error) {
	var data []models.Attendance
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Attendance{}).
		Where("user_id = ?", userId).
		Preload("User")

	dataAkhir := response.PaginateResponseDto[[]models.Attendance]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "attendance_id"
		queryParams.SortBy = &sort
	}

	allowedDynamicList := attendanceDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Where("user_id = ?", userId).
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
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
		"location_code": {
			Field: "location_code",
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
