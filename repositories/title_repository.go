package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type TitleRepository interface {
	Create(data *models.Title) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Title], error)
	Update(data *models.Title) error
	Delete(id string) error
}

type titleRepository struct {
	db *gorm.DB
}

func NewTitleRepository(db *gorm.DB) TitleRepository {
	return &titleRepository{db}
}

func (r *titleRepository) Create(data *models.Title) error {
	return r.db.Create(data).Error
}

func (r *titleRepository) Update(data *models.Title) error {
	return r.db.Model(&models.Title{}).Where("title_id = ?", data.TitleId).Updates(data).Error
}

func (r *titleRepository) Delete(id string) error {
	return r.db.Delete(&models.Title{}, "title_id = ?", id).Error
}

func (r *titleRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Title], error) {
	var data []models.Title
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.Title{})
	result := response.PaginateResponseDto[[]models.Title]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "title_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			title_id::text LIKE ? OR company_code LIKE ? OR branch_code LIKE ? OR
			office_code LIKE ? OR division_code LIKE ? OR department_code LIKE ? OR
			title_code LIKE ? OR title_name LIKE ? OR object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"title_id":        {Field: "title_id", Query: " = ?"},
		"company_code":    {Field: "company_code", Query: " LIKE ?"},
		"branch_code":     {Field: "branch_code", Query: " LIKE ?"},
		"office_code":     {Field: "office_code", Query: " LIKE ?"},
		"division_code":   {Field: "division_code", Query: " LIKE ?"},
		"department_code": {Field: "department_code", Query: " LIKE ?"},
		"title_code":      {Field: "title_code", Query: " LIKE ?"},
		"title_name":      {Field: "title_name", Query: " LIKE ?"},
		"object_code":     {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":    {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
