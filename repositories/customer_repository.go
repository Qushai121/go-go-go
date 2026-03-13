package repositories

import (
	"hrms_go/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindAll() ([]models.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) Create(customer *models.Customer) error {
	return c.db.Create(customer).Error
}

func (c *customerRepository) FindAll() ([]models.Customer, error) {
	var customers []models.Customer
	err := c.db.Find(&customers).Error
	return customers, err
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}