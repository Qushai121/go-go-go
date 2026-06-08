package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type DivisionRepository interface {
	Create(data *models.Division) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Division], error)
	Update(data *models.Division) error
	Delete(id string) error
}

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) DivisionRepository {
	return &divisionRepository{db}
}

func (r *divisionRepository) Create(data *models.Division) error {
	return r.db.Create(data).Error
}

func (r *divisionRepository) Update(data *models.Division) error {
	return r.db.Model(&models.Division{}).Where("division_id = ?", data.DivisionId).Updates(data).Error
}

func (r *divisionRepository) Delete(id string) error {
	return r.db.Delete(&models.Division{}, "division_id = ?", id).Error
}

func (r *divisionRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Division], error) {
	var data []models.Division
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.Division{})
	result := response.PaginateResponseDto[[]models.Division]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "division_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			division_id::text ILIKE ? OR company_code ILIKE ? OR
			division_code ILIKE ? OR division_name ILIKE ? OR
			object_code ILIKE ? OR timezone_set ILIKE ?
		`, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"division_id":   {Field: "division_id", Query: " = ?"},
		"company_code":  {Field: "company_code", Query: " ILIKE ?"},
		"division_code": {Field: "division_code", Query: " ILIKE ?"},
		"division_name": {Field: "division_name", Query: " ILIKE ?"},
		"object_code":   {Field: "object_code", Query: " ILIKE ?"},
		"timezone_set":  {Field: "timezone_set", Query: " ILIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
