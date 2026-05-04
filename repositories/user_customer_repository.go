package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserCustomerRepository interface {
	Create(data *models.UserCustomer) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserCustomer], error)
	Update(data *models.UserCustomer) error
	Delete(id string) error
}

type userCustomerRepository struct {
	db *gorm.DB
}

func NewUserCustomerRepository(db *gorm.DB) UserCustomerRepository {
	return &userCustomerRepository{db}
}

func (r *userCustomerRepository) Create(data *models.UserCustomer) error {
	return r.db.Create(data).Error
}

func (r *userCustomerRepository) Update(data *models.UserCustomer) error {
	return r.db.Model(&models.UserCustomer{}).Where("user_customer_id = ?", data.UserCustomerId).Updates(data).Error
}

func (r *userCustomerRepository) Delete(id string) error {
	return r.db.Delete(&models.UserCustomer{}, "user_customer_id = ?", id).Error
}

func (r *userCustomerRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserCustomer], error) {
	var data []models.UserCustomer
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.UserCustomer{})
	result := response.PaginateResponseDto[[]models.UserCustomer]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "user_customer_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			user_customer_id::text LIKE ? OR employee_nik LIKE ? OR customer_code LIKE ? OR
			object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"user_customer_id": {Field: "user_customer_id", Query: " = ?"},
		"employee_nik":     {Field: "employee_nik", Query: " LIKE ?"},
		"customer_code":    {Field: "customer_code", Query: " LIKE ?"},
		"object_code":      {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":     {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
