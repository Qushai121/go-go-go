package repositories

import (
	"strings"

	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserSiteRepository interface {
	Create(data *models.UserSite) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserSite], error)
	Update(data *models.UserSite) error
	Delete(id string) error
}

type userSiteRepository struct {
	db *gorm.DB
}

func NewUserSiteRepository(db *gorm.DB) UserSiteRepository {
	return &userSiteRepository{db}
}

func (r *userSiteRepository) Create(data *models.UserSite) error {
	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if data.ObjectCode == "" {
		data.ObjectCode = "USER_SITE"
	}
	return r.db.Create(data).Error
}

func (r *userSiteRepository) Update(data *models.UserSite) error {
	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	return r.db.Model(&models.UserSite{}).Where("user_site_id = ?", data.UserSiteId).Updates(data).Error
}

func (r *userSiteRepository) Delete(id string) error {
	return r.db.Delete(&models.UserSite{}, "user_site_id = ?", id).Error
}

func (r *userSiteRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserSite], error) {
	var data []models.UserSite
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Table("hrms_user_site us").
		Select(`
			us.*,
			s.site_id AS site_id,
			s.site_name AS site_name,
			s.site_phone AS site_phone,
			s.site_address AS site_address,
			s.site_latitude AS site_latitude,
			s.site_longitude AS site_longitude,
			s.max_radius AS max_radius
		`).
		Joins(`
			LEFT JOIN hrms_site s
				ON s.company_code = us.company_code
				AND s.site_type = us.site_type
				AND s.site_code = us.site_code
		`)

	result := response.PaginateResponseDto[[]models.UserSite]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "us.user_site_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			us.user_site_id::text ILIKE ? OR
			us.company_code ILIKE ? OR
			us.employee_nik ILIKE ? OR
			us.site_type ILIKE ? OR
			us.site_code ILIKE ? OR
			s.site_name ILIKE ? OR
			us.object_code ILIKE ? OR
			us.timezone_set ILIKE ?
		`, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"user_site_id": {Field: "us.user_site_id", Query: " = ?"},
		"company_code": {Field: "us.company_code", Query: " ILIKE ?"},
		"employee_nik": {Field: "us.employee_nik", Query: " ILIKE ?"},
		"site_type":    {Field: "us.site_type", Query: " ILIKE ?"},
		"site_code":    {Field: "us.site_code", Query: " ILIKE ?"},
		"site_name":    {Field: "s.site_name", Query: " ILIKE ?"},
		"object_code":  {Field: "us.object_code", Query: " ILIKE ?"},
		"timezone_set": {Field: "us.timezone_set", Query: " ILIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
