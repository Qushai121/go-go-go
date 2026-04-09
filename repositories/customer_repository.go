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
	return c.db.Model(&models.Customer{}).Where("customer_id = ?",customer.CustomerId).Updates(&customer).Error
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
		Data:        	data,
		TotalRecord: 	totalRecord,
		TotalPage: 		totalPage,
	}
	
	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}
