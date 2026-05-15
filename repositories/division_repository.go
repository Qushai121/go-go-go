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
			division_id::text LIKE ? OR company_code LIKE ? OR branch_code LIKE ? OR
			office_code LIKE ? OR division_code LIKE ? OR division_name LIKE ? OR
			object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"division_id":   {Field: "division_id", Query: " = ?"},
		"company_code":  {Field: "company_code", Query: " LIKE ?"},
		"branch_code":   {Field: "branch_code", Query: " LIKE ?"},
		"office_code":   {Field: "office_code", Query: " LIKE ?"},
		"division_code": {Field: "division_code", Query: " LIKE ?"},
		"division_name": {Field: "division_name", Query: " LIKE ?"},
		"object_code":   {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":  {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
