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
	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			title_id::text ILIKE ? OR company_code ILIKE ? OR
			division_code ILIKE ? OR department_code ILIKE ? OR
			title_code ILIKE ? OR title_name ILIKE ? OR object_code ILIKE ? OR timezone_set ILIKE ?
		`, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"title_id":        {Field: "title_id", Query: " = ?"},
		"company_code":    {Field: "company_code", Query: " ILIKE ?"},
		"division_code":   {Field: "division_code", Query: " ILIKE ?"},
		"department_code": {Field: "department_code", Query: " ILIKE ?"},
		"title_code":      {Field: "title_code", Query: " ILIKE ?"},
		"title_name":      {Field: "title_name", Query: " ILIKE ?"},
		"object_code":     {Field: "object_code", Query: " ILIKE ?"},
		"timezone_set":    {Field: "timezone_set", Query: " ILIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
