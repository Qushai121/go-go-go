package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserOfficeRepository interface {
	Create(data *models.UserOffice) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserOffice], error)
	Update(data *models.UserOffice) error
	Delete(id string) error
}

type userOfficeRepository struct {
	db *gorm.DB
}

func NewUserOfficeRepository(db *gorm.DB) UserOfficeRepository {
	return &userOfficeRepository{db}
}

func (r *userOfficeRepository) Create(data *models.UserOffice) error {
	return r.db.Create(data).Error
}

func (r *userOfficeRepository) Update(data *models.UserOffice) error {
	return r.db.Model(&models.UserOffice{}).Where("user_office_id = ?", data.UserOfficeId).Updates(data).Error
}

func (r *userOfficeRepository) Delete(id string) error {
	return r.db.Delete(&models.UserOffice{}, "user_office_id = ?", id).Error
}

func (r *userOfficeRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserOffice], error) {
	var data []models.UserOffice
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.UserOffice{})
	result := response.PaginateResponseDto[[]models.UserOffice]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "user_office_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			user_office_id::text LIKE ? OR company_code LIKE ? OR branch_code LIKE ? OR
			employee_nik LIKE ? OR office_code LIKE ? OR object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"user_office_id": {Field: "user_office_id", Query: " = ?"},
		"company_code":   {Field: "company_code", Query: " LIKE ?"},
		"branch_code":    {Field: "branch_code", Query: " LIKE ?"},
		"employee_nik":   {Field: "employee_nik", Query: " LIKE ?"},
		"office_code":    {Field: "office_code", Query: " LIKE ?"},
		"object_code":    {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":   {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
