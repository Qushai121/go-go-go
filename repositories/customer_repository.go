package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Customer], error)
	Delete(customerId string) error
	Update(customer *models.Customer) error
}

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) Delete(customerId string) error {
	return c.db.Delete(&models.Customer{}, "customer_id = ?", customerId).Error
}

func (c *customerRepository) Update(customer *models.Customer) error {
	return c.db.Model(&models.Customer{}).Where("customer_id = ?", customer.CustomerId).Updates(customer).Error
}

func (c *customerRepository) Create(customer *models.Customer) error {
	return c.db.Create(customer).Error
}

func (c *customerRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Customer], error) {
	var data []models.Customer
	var totalRecord int64
	var totalPage int

	modelDb := c.db.Model(&models.Customer{})

	if queryParams.SortBy == nil {
		sort := "customer_id"
		queryParams.SortBy = &sort
	}

	result := response.PaginateResponseDto[[]models.Customer]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			customer_id::text LIKE ? OR
			customer_code LIKE ? OR
			customer_name LIKE ? OR
			customer_address LIKE ? OR
			customer_latitude LIKE ? OR
			customer_longitude LIKE ? OR
			max_radius::text LIKE ? OR
			object_code LIKE ? OR
			timezone_set LIKE ?
		`, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"customer_id": {
			Field: "customer_id",
			Query: " = ?",
		},
		"customer_code": {
			Field: "customer_code",
			Query: " LIKE ?",
		},
		"customer_name": {
			Field: "customer_name",
			Query: " LIKE ?",
		},
		"customer_address": {
			Field: "customer_address",
			Query: " LIKE ?",
		},
		"customer_latitude": {
			Field: "customer_latitude",
			Query: " LIKE ?",
		},
		"customer_longitude": {
			Field: "customer_longitude",
			Query: " LIKE ?",
		},
		"max_radius": {
			Field: "max_radius",
			Query: " = ?",
		},
		"object_code": {
			Field: "object_code",
			Query: " LIKE ?",
		},
		"timezone_set": {
			Field: "timezone_set",
			Query: " LIKE ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}
