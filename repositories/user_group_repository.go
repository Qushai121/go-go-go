package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserGroupRepository interface {
	Create(data *models.UserGroup) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserGroup], error)
	Update(data *models.UserGroup) error
	Delete(id string) error
}

type userGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(db *gorm.DB) UserGroupRepository {
	return &userGroupRepository{db}
}

func (r *userGroupRepository) Create(data *models.UserGroup) error {
	return r.db.Create(data).Error
}

func (r *userGroupRepository) Update(data *models.UserGroup) error {
	return r.db.Model(&models.UserGroup{}).Where("usergroup_id = ?", data.UserGroupId).Updates(data).Error
}

func (r *userGroupRepository) Delete(id string) error {
	return r.db.Delete(&models.UserGroup{}, "usergroup_id = ?", id).Error
}

func (r *userGroupRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.UserGroup], error) {
	var data []models.UserGroup
	var totalRecord int64
	var totalPage int
	modelDb := r.db.Model(&models.UserGroup{})
	result := response.PaginateResponseDto[[]models.UserGroup]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "usergroup_id"
		queryParams.SortBy = &sort
	}
	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			usergroup_id::text LIKE ? OR company_code LIKE ? OR usergroup_code LIKE ? OR
			usergroup_name LIKE ? OR object_code LIKE ? OR timezone_set LIKE ?
		`, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"usergroup_id":   {Field: "usergroup_id", Query: " = ?"},
		"company_code":   {Field: "company_code", Query: " LIKE ?"},
		"usergroup_code": {Field: "usergroup_code", Query: " LIKE ?"},
		"usergroup_name": {Field: "usergroup_name", Query: " LIKE ?"},
		"object_code":    {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":   {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
