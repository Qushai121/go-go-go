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
		Where("user_id = ? AND date = ? AND check_type = ?",
			attendance.UserId,
			attendance.Date,
			attendance.CheckType).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("attendance already exists for today")
	}

	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Attendance], error) {
	var data []models.Attendance
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Attendance{}).Preload("User")

	dataAkhir := response.PaginateResponseDto[[]models.Attendance]{
		Data:        	data,
		TotalRecord: 	totalRecord,
		TotalPage: 		totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "attendance_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
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
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "attendance_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Where("user_id = ?",userId).Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db}
}
