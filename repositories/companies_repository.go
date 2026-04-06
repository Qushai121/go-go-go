package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type CompaniesRepository interface {
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Companies], error)
	Create(companies *models.Companies) error
	Delete(companiesId string) error
	Update(companies *models.Companies) error
}

type companiesRepository struct {
	db *gorm.DB
}

func (c *companiesRepository) Create(companies *models.Companies) error {
	return c.db.Create(companies).Error
}

func (c *companiesRepository) Delete(companiesId string) error {
	return c.db.Delete(&models.Companies{}, "companies_id = ?", companiesId).Error
}

func (c *companiesRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Companies], error) {
	data := []models.Companies{}
	modelDb := c.db.Model(&models.Companies{})
	var totalRecord int64

	dataAkhir := response.PaginateResponseDto[[]models.Companies]{
		Data:        data,
		TotalRecord: totalRecord,
	}

	if queryParams.SortBy == nil {
		sort := "companies_id"
		queryParams.SortBy = &sort
	}

	if err := utils.GetQuery(queryParams, modelDb, &totalRecord).Find(&data).Error; err != nil {
		return dataAkhir, err
	}

	return dataAkhir, nil
}

func (c *companiesRepository) Update(companies *models.Companies) error {
	return c.db.Model(&models.Companies{}).Updates(&companies).Error
}

func NewCompaniesRepository(db *gorm.DB) CompaniesRepository {
	return &companiesRepository{db}
}
