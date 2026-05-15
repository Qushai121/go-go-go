package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(data *models.Branch) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Branch], error)
	Update(data *models.Branch) error
	Delete(id string) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db}
}

func (r *branchRepository) Create(data *models.Branch) error {
	return r.db.Create(data).Error
}

func (r *branchRepository) Update(data *models.Branch) error {
	return r.db.Model(&models.Branch{}).Where("branch_id = ?", data.BranchId).Updates(data).Error
}

func (r *branchRepository) Delete(id string) error {
	return r.db.Delete(&models.Branch{}, "branch_id = ?", id).Error
}

func (r *branchRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Branch], error) {
	var data []models.Branch
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Branch{})
	result := response.PaginateResponseDto[[]models.Branch]{Data: data, TotalRecord: totalRecord, TotalPage: totalPage}

	if queryParams.SortBy == nil {
		sort := "branch_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			branch_id::text LIKE ? OR
			company_code LIKE ? OR
			branch_code LIKE ? OR
			branch_name LIKE ? OR
			branch_address LIKE ? OR
			object_code LIKE ? OR
			timezone_set LIKE ?
		`, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"branch_id":      {Field: "branch_id", Query: " = ?"},
		"company_code":   {Field: "company_code", Query: " LIKE ?"},
		"branch_code":    {Field: "branch_code", Query: " LIKE ?"},
		"branch_name":    {Field: "branch_name", Query: " LIKE ?"},
		"branch_address": {Field: "branch_address", Query: " LIKE ?"},
		"object_code":    {Field: "object_code", Query: " LIKE ?"},
		"timezone_set":   {Field: "timezone_set", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}
