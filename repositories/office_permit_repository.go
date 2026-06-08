package repositories

import (
	"fmt"
	"strings"

	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type OfficePermitRepository interface {
	Create(data *models.OfficePermit) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.OfficePermit], error)
	Update(data *models.OfficePermit) error
	Delete(id string) error
}

type officePermitRepository struct {
	db *gorm.DB
}

func NewOfficePermitRepository(db *gorm.DB) OfficePermitRepository {
	return &officePermitRepository{db}
}

func (r *officePermitRepository) Create(data *models.OfficePermit) error {
	data.EmployeeNik = strings.TrimSpace(data.EmployeeNik)
	data.Status = strings.ToUpper(strings.TrimSpace(data.Status))
	data.ObjectCode = strings.TrimSpace(data.ObjectCode)
	data.CreatedBy = strings.TrimSpace(data.CreatedBy)

	if data.EmployeeNik == "" {
		return fmt.Errorf("employee_nik wajib diisi")
	}
	if data.OfficePermitDate.IsZero() {
		return fmt.Errorf("office_permit_date wajib diisi")
	}
	if data.Status == "" {
		data.Status = "P"
	}
	if data.CurrentStep <= 0 {
		data.CurrentStep = 1
	}
	if data.ObjectCode == "" {
		data.ObjectCode = "LEAVE_HISTORY"
	}
	if data.CreatedBy == "" {
		data.CreatedBy = "System"
	}

	return r.db.Create(data).Error
}

func (r *officePermitRepository) Update(data *models.OfficePermit) error {
	data.EmployeeNik = strings.TrimSpace(data.EmployeeNik)
	data.Status = strings.ToUpper(strings.TrimSpace(data.Status))
	data.ObjectCode = strings.TrimSpace(data.ObjectCode)

	return r.db.Model(&models.OfficePermit{}).
		Where("office_permit_id = ?", data.OfficePermitId).
		Updates(data).Error
}

func (r *officePermitRepository) Delete(id string) error {
	return r.db.Delete(&models.OfficePermit{}, "office_permit_id = ?", id).Error
}

func (r *officePermitRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.OfficePermit], error) {
	var data []models.OfficePermit
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.OfficePermit{})
	result := response.PaginateResponseDto[[]models.OfficePermit]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "office_permit_date"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			office_permit_id::text ILIKE ? OR
			employee_nik ILIKE ? OR
			office_permit_date::text ILIKE ? OR
			remarks ILIKE ? OR
			status ILIKE ? OR
			current_step::text ILIKE ? OR
			approvalheader_id::text ILIKE ? OR
			object_code ILIKE ? OR
			created_by ILIKE ? OR
			created_at::text ILIKE ? OR
			updated_by ILIKE ? OR
			updated_at::text ILIKE ?
		`, search, search, search, search, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"office_permit_id":   {Field: "office_permit_id", Query: " = ?"},
		"employee_nik":       {Field: "employee_nik", Query: " ILIKE ?"},
		"office_permit_date": {Field: "office_permit_date::date", Query: " = CAST(? AS date)"},
		"date_start":         {Field: "office_permit_date::date", Query: " >= CAST(? AS date)"},
		"date_end":           {Field: "office_permit_date::date", Query: " <= CAST(? AS date)"},
		"remarks":            {Field: "remarks", Query: " ILIKE ?"},
		"status":             {Field: "status", Query: " ILIKE ?"},
		"current_step":       {Field: "current_step", Query: " = ?"},
		"approvalheader_id":  {Field: "approvalheader_id", Query: " = ?"},
		"object_code":        {Field: "object_code", Query: " ILIKE ?"},
		"created_by":         {Field: "created_by", Query: " ILIKE ?"},
		"updated_by":         {Field: "updated_by", Query: " ILIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).
		Find(&result.Data).Error

	return result, err
}
