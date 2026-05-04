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
	var data = []models.Companies{}
	var totalRecord int64
	var totalPage int
	modelDb := c.db.Model(&models.Companies{})

	dataAkhir := response.PaginateResponseDto[[]models.Companies]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "companies_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			companies_id::text LIKE ? OR
			companies_code LIKE ? OR
			companies_name LIKE ?
		`, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"companies_id": {
			Field: "companies_id",
			Query: " = ?",
		},
		"companies_code": {
			Field: "companies_code",
			Query: " LIKE ?",
		},
		"companies_name": {
			Field: "companies_name",
			Query: " LIKE ?",
		},
	}

	if err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error; err != nil {
		return dataAkhir, err
	}

	return dataAkhir, nil
}

func (c *companiesRepository) Update(companies *models.Companies) error {
	return c.db.Model(&models.Companies{}).Where("companies_id = ?", companies.CompaniesId).Updates(&companies).Error
}

func NewCompaniesRepository(db *gorm.DB) CompaniesRepository {
	return &companiesRepository{db}
}
