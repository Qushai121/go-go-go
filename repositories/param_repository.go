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
	var totalPage int

	modelDb := r.db.Model(&models.Param{}).Preload("ParamGroup")

	dataAkhir := response.PaginateResponseDto[[]models.Param]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "param_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			param_id::text LIKE ? OR
			paramgroup_id::text LIKE ? OR
			param_code LIKE ? OR
			param_name LIKE ?
		`, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"param_id": {
			Field: "param_id",
			Query: " = ?",
		},
		"paramgroup_id": {
			Field: "paramgroup_id",
			Query: " = ?",
		},
		"param_code": {
			Field: "param_code",
			Query: " LIKE ?",
		},
		"param_name": {
			Field: "param_name",
			Query: " LIKE ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func NewParamRepository(db *gorm.DB) ParamRepository {
	return &paramRepository{db}
}
