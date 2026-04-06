package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ParamRepository interface {
	Create(param *models.Param) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Param], error)
	Update(param *models.Param) error
	Delete(paramId string) error
}

type paramRepository struct {
	db *gorm.DB
}

func (r *paramRepository) Create(param *models.Param) error {
	return r.db.Create(param).Error
}

func (r *paramRepository) Update(param *models.Param) error {
	return r.db.Model(&models.Param{}).
		Where("param_id = ?", param.ParamId).
		Updates(param).Error
}

func (r *paramRepository) Delete(paramId string) error {
	return r.db.Delete(&models.Param{}, "param_id = ?", paramId).Error
}

func (r *paramRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Param], error) {
	var data []models.Param
	var totalRecord int64

	modelDb := r.db.Model(&models.Param{}).Preload("ParamGroup")

	dataAkhir := response.PaginateResponseDto[[]models.Param]{
		Data:        data,
		TotalRecord: totalRecord,
	}

	if queryParams.SortBy == nil {
		sort := "param_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &totalRecord).Find(&data).Error

	dataAkhir.Data = data
	dataAkhir.TotalRecord = totalRecord

	return dataAkhir, err
}

func NewParamRepository(db *gorm.DB) ParamRepository {
	return &paramRepository{db}
}
