package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserCompanyRepository interface {
	Create(data *models.UserCompany) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserCompany], error)
	Update(data *models.UserCompany) error
	Delete(id string) error
}

type userCompanyRepository struct {
	db *gorm.DB
}

func NewUserCompanyRepository(db *gorm.DB) UserCompanyRepository {
	return &userCompanyRepository{db}
}

func (r *userCompanyRepository) Create(data *models.UserCompany) error {
	return r.db.Create(data).Error
}

func (r *userCompanyRepository) Update(data *models.UserCompany) error {
	return r.db.Model(&models.UserCompany{}).Where("user_company_id = ?", data.UserCompanyId).Updates(data).Error
}

func (r *userCompanyRepository) Delete(id string) error {
	return r.db.Delete(&models.UserCompany{}, "user_company_id = ?", id).Error
}

func (r *userCompanyRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserCompany], error) {
	var data []models.UserCompany
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.UserCompany{})
	result := response.PaginateResponseDto[[]models.UserCompany]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "user_company_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			user_company_id::text LIKE ? OR employee_nik LIKE ? OR company_code LIKE ? OR
			object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"user_company_id": {Field: "user_company_id", Query: " = ?"},
		"employee_nik":    {Field: "employee_nik", Query: " LIKE ?"},
		"company_code":    {Field: "company_code", Query: " LIKE ?"},
		"object_code":     {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":    {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
