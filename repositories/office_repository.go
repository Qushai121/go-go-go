package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type OfficeRepository interface {
	Create(office *models.Office) error
	Update(office *models.Office) error
	Delete(id string) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Office], error)
}

type officeRepository struct {
	db *gorm.DB
}

func (r *officeRepository) Create(office *models.Office) error {
	return r.db.Create(office).Error
}

func (r *officeRepository) Update(office *models.Office) error {
	return r.db.Model(&models.Office{}).
		Where("office_id = ?", office.OfficeId).
		Updates(office).Error
}

func (r *officeRepository) Delete(id string) error {
	return r.db.Delete(&models.Office{}, "office_id = ?", id).Error
}

func (r *officeRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Office], error) {
	var data []models.Office
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Office{})

	dataAkhir := response.PaginateResponseDto[[]models.Office]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "office_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			office_id::text LIKE ? OR
			company_code LIKE ? OR
			branch_code LIKE ? OR
			office_code LIKE ? OR
			office_name LIKE ? OR
			office_phone LIKE ? OR
			office_address LIKE ? OR
			office_latitude LIKE ? OR
			office_longitude LIKE ? OR
			object_code LIKE ? OR
			timezone_set LIKE ?
		`,
			search, search, search, search, search,
			search, search, search, search, search,
			search, search,
		)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"office_id": {
			Field: "office_id",
			Query: " = ?",
		},
		"company_code": {
			Field: "company_code",
			Query: " LIKE ?",
		},
		"branch_code": {
			Field: "branch_code",
			Query: " LIKE ?",
		},
		"office_code": {
			Field: "office_code",
			Query: " LIKE ?",
		},
		"office_name": {
			Field: "office_name",
			Query: " LIKE ?",
		},
		"office_phone": {
			Field: "office_phone",
			Query: " LIKE ?",
		},
		"office_address": {
			Field: "office_address",
			Query: " LIKE ?",
		},
		"office_latitude": {
			Field: "office_latitude",
			Query: " LIKE ?",
		},
		"office_longitude": {
			Field: "office_longitude",
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

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewOfficeRepository(db *gorm.DB) OfficeRepository {
	return &officeRepository{db}
}
