package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ApprovalTemplateRepository interface {
	CreateHeader(data *models.ApprovalTemplateHeader) error
	UpdateHeader(id string, data *models.ApprovalTemplateHeader) error
	DeleteHeader(id string) error
	FindAllHeader(query *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalTemplateHeader], error)
	DetailHeader(id string) (models.ApprovalTemplateHeader, error)

	CreateDetail(data *models.ApprovalTemplateDetail) error
	UpdateDetail(id string, data *models.ApprovalTemplateDetail) error
	DeleteDetail(id string) error
	FindByHeader(headerId string) ([]models.ApprovalTemplateDetail, error)
}

type approvalTemplateRepository struct {
	db *gorm.DB
}

func (r *approvalTemplateRepository) CreateHeader(data *models.ApprovalTemplateHeader) error {
	return r.db.Create(data).Error
}

func (r *approvalTemplateRepository) UpdateHeader(id string, data *models.ApprovalTemplateHeader) error {
	return r.db.Model(&models.ApprovalTemplateHeader{}).
		Where("approval_template_header_id = ?", id).
		Updates(data).Error
}

func (r *approvalTemplateRepository) DeleteHeader(id string) error {
	return r.db.Delete(&models.ApprovalTemplateHeader{}, "approval_template_header_id = ?", id).Error
}

func (r *approvalTemplateRepository) DetailHeader(id string) (models.ApprovalTemplateHeader, error) {
	var data models.ApprovalTemplateHeader
	err := r.db.
		Where("approval_template_header_id = ?", id).
		Find(&data).Error
	return data, err
}

func (r *approvalTemplateRepository) FindAllHeader(query *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalTemplateHeader], error) {
	var data []models.ApprovalTemplateHeader
	var total int64
	var totalPage int

	db := r.db.Model(&models.ApprovalTemplateHeader{})

	if query.Search != nil && *query.Search != "" {
		search := "%" + *query.Search + "%"
		db = db.Where(`
			approval_template_header_id::text LIKE ? OR
			template_type LIKE ? OR
			created_by LIKE ? OR
			updated_by LIKE ?
		`, search, search, search, search)
	}

	result :=  response.PaginateResponseDto[[]models.ApprovalTemplateHeader]{
		Data:        data,
		TotalRecord: total,
		TotalPage:   totalPage,
	}

	err := utils.GetQuery(query, db, &result.TotalRecord, &result.TotalPage).
		Find(&result.Data).Error

	return result,err
}

func (r *approvalTemplateRepository) CreateDetail(data *models.ApprovalTemplateDetail) error {
	return r.db.Create(data).Error
}

func (r *approvalTemplateRepository) UpdateDetail(id string, data *models.ApprovalTemplateDetail) error {
	return r.db.Model(&models.ApprovalTemplateDetail{}).
		Where("approval_template_detail_id = ?", id).
		Updates(data).Error
}

func (r *approvalTemplateRepository) DeleteDetail(id string) error {
	return r.db.Delete(&models.ApprovalTemplateDetail{}, "approval_template_detail_id = ?", id).Error
}

func (r *approvalTemplateRepository) FindByHeader(headerId string) ([]models.ApprovalTemplateDetail, error) {
	var data []models.ApprovalTemplateDetail

	err := r.db.
		Preload("Approver").
		Where("approval_template_header_id = ?", headerId).
		Order("sequence_number ASC").
		Find(&data).Error

	return data, err
}

func NewApprovalTemplateRepository(db *gorm.DB) ApprovalTemplateRepository {
	return &approvalTemplateRepository{db}
}