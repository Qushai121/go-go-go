package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ParamGroupRepository interface {
	Create(paramGroup *models.ParamGroup) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ParamGroup], error)
	Delete(paramGroupId string) error
	Update(paramGroup *models.ParamGroup) error
}

type paramGroupRepository struct {
	db *gorm.DB
}

func (r *paramGroupRepository) Create(paramGroup *models.ParamGroup) error {
	return r.db.Create(paramGroup).Error
}

func (r *paramGroupRepository) Delete(paramGroupId string) error {
	return r.db.Delete(&models.ParamGroup{}, "paramgroup_id = ?", paramGroupId).Error
}

func (r *paramGroupRepository) Update(paramGroup *models.ParamGroup) error {
	return r.db.Model(&models.ParamGroup{}).
		Where("paramgroup_id = ?", paramGroup.ParamGroupId).
		Updates(paramGroup).Error
}

func (r *paramGroupRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ParamGroup], error) {
	var data []models.ParamGroup
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.ParamGroup{})

	dataAkhir := response.PaginateResponseDto[[]models.ParamGroup]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "paramgroup_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			paramgroup_id::text LIKE ? OR
			company_code LIKE ? OR
			paramgroup_code LIKE ? OR
			paramgroup_name LIKE ?
		`, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"paramgroup_id": {
			Field: "paramgroup_id",
			Query: " = ?",
		},
		"company_code": {
			Field: "company_code",
			Query: " LIKE ?",
		},
		"paramgroup_code": {
			Field: "paramgroup_code",
			Query: " LIKE ?",
		},
		"paramgroup_name": {
			Field: "paramgroup_name",
			Query: " LIKE ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func NewParamGroupRepository(db *gorm.DB) ParamGroupRepository {
	return &paramGroupRepository{db}
}
