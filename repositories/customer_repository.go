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
	return c.db.Model(&models.Customer{}).Where("customer_id = ?", customer.CustomerId).Updates(&customer).Error
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

	dataAkhir := response.PaginateResponseDto[[]models.Customer]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			customer_id::text LIKE ? OR
			location_code LIKE ? OR
			location_name LIKE ? OR
			address LIKE ? OR
			target_latitude LIKE ? OR
			target_longitude LIKE ? OR
			radius_meter::text LIKE ?
		`, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"customer_id": {
			Field: "customer_id",
			Query: " = ?",
		},
		"location_code": {
			Field: "location_code",
			Query: " LIKE ?",
		},
		"location_name": {
			Field: "location_name",
			Query: " LIKE ?",
		},
		"address": {
			Field: "address",
			Query: " LIKE ?",
		},
		"target_latitude": {
			Field: "target_latitude",
			Query: " LIKE ?",
		},
		"target_longitude": {
			Field: "target_longitude",
			Query: " LIKE ?",
		},
		"radius_meter": {
			Field: "radius_meter",
			Query: " = ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}
