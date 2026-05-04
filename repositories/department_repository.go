package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(data *models.Department) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Department], error)
	Update(data *models.Department) error
	Delete(id string) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db}
}

func (r *departmentRepository) Create(data *models.Department) error {
	return r.db.Create(data).Error
}

func (r *departmentRepository) Update(data *models.Department) error {
	return r.db.Model(&models.Department{}).Where("department_id = ?", data.DepartmentId).Updates(data).Error
}

func (r *departmentRepository) Delete(id string) error {
	return r.db.Delete(&models.Department{}, "department_id = ?", id).Error
}

func (r *departmentRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Department], error) {
	var data []models.Department
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.Department{})
	result := response.PaginateResponseDto[[]models.Department]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "department_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			department_id::text LIKE ? OR company_code LIKE ? OR branch_code LIKE ? OR
			office_code LIKE ? OR division_code LIKE ? OR department_code LIKE ? OR
			department_name LIKE ? OR object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"department_id":   {Field: "department_id", Query: " = ?"},
		"company_code":    {Field: "company_code", Query: " LIKE ?"},
		"branch_code":     {Field: "branch_code", Query: " LIKE ?"},
		"office_code":     {Field: "office_code", Query: " LIKE ?"},
		"division_code":   {Field: "division_code", Query: " LIKE ?"},
		"department_code": {Field: "department_code", Query: " LIKE ?"},
		"department_name": {Field: "department_name", Query: " LIKE ?"},
		"object_code":     {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":    {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
