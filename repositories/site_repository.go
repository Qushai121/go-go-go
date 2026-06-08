package repositories

import (
	"strings"

	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type SiteRepository interface {
	Create(data *models.Site) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Site], error)
	Update(data *models.Site) error
	Delete(id string) error
}

type siteRepository struct {
	db *gorm.DB
}

func NewSiteRepository(db *gorm.DB) SiteRepository {
	return &siteRepository{db}
}

func (r *siteRepository) Create(data *models.Site) error {
	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if data.ObjectCode == "" {
		data.ObjectCode = "SITE"
	}
	if data.MaxRadius <= 0 {
		data.MaxRadius = 5
	}
	return r.db.Create(data).Error
}

func (r *siteRepository) Update(data *models.Site) error {
	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	return r.db.Model(&models.Site{}).Where("site_id = ?", data.SiteId).Updates(data).Error
}

func (r *siteRepository) Delete(id string) error {
	return r.db.Delete(&models.Site{}, "site_id = ?", id).Error
}

func (r *siteRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Site], error) {
	var data []models.Site
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Site{})
	result := response.PaginateResponseDto[[]models.Site]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "site_name"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			site_id::text ILIKE ? OR
			company_code ILIKE ? OR
			site_type ILIKE ? OR
			site_code ILIKE ? OR
			site_name ILIKE ? OR
			site_phone ILIKE ? OR
			site_address ILIKE ? OR
			site_latitude ILIKE ? OR
			site_longitude ILIKE ? OR
			object_code ILIKE ? OR
			timezone_set ILIKE ?
		`, search, search, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"site_id":        {Field: "site_id", Query: " = ?"},
		"company_code":   {Field: "company_code", Query: " ILIKE ?"},
		"site_type":      {Field: "site_type", Query: " ILIKE ?"},
		"site_code":      {Field: "site_code", Query: " ILIKE ?"},
		"site_name":      {Field: "site_name", Query: " ILIKE ?"},
		"site_phone":     {Field: "site_phone", Query: " ILIKE ?"},
		"site_address":   {Field: "site_address", Query: " ILIKE ?"},
		"site_latitude":  {Field: "site_latitude", Query: " ILIKE ?"},
		"site_longitude": {Field: "site_longitude", Query: " ILIKE ?"},
		"object_code":    {Field: "object_code", Query: " ILIKE ?"},
		"timezone_set":   {Field: "timezone_set", Query: " ILIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
