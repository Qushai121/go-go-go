package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type LeaveRepository interface {
	Create(leave *models.Leave) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Leave], error)
	Delete(leaveId string) error
	Update(leave *models.Leave) error
}

type leaveRepository struct {
	db *gorm.DB
}

func (r *leaveRepository) Create(leave *models.Leave) error {
	return r.db.Create(leave).Error
}

func (r *leaveRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Leave], error) {
	var data []models.Leave
	modelDb := r.db.Model(&models.Leave{})
	var totalRecord int64

	dataAkhir := response.PaginateResponseDto[[]models.Leave]{
		Data:        data,
		TotalRecord: totalRecord,
	}

	if queryParams.SortBy == nil {
		sort := "leave_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &totalRecord).Find(&data).Error
	return dataAkhir, err
}

func (r *leaveRepository) Update(leave *models.Leave) error {
	return r.db.Model(&models.Leave{}).Updates(&leave).Error
}

func (r *leaveRepository) Delete(leaveId string) error {
	return r.db.Delete(&models.Leave{}, "leave_id = ?", leaveId).Error
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{db}
}