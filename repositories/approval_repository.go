package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/approval"
	"hrms_go/dto/response"
	mappers "hrms_go/mapper"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ApprovalRepository interface {
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalHeader], error)
	Detail(approvalHeaderId string) (response.PaginateResponseDto[[]approval.ApprovalDetailResponseDto], error)
	Delete(approvalHeaderId string) error
	Approve(approvalDetail models.ApprovalDetail) error
}

type approvalRepository struct {
	db *gorm.DB
}

func (a *approvalRepository) Approve(approvalDetail models.ApprovalDetail) error {
	approvalHeader := models.ApprovalHeader{}

	if err := a.db.Model(&models.ApprovalHeader{}).
		Where("approval_header_id = ?", approvalDetail.ApprovalHeaderId).
		First(&approvalHeader).Error; err != nil {
		return err
	}

	approvalTemplateDetailData := models.ApprovalTemplateDetail{}
	if err := a.db.Model(&models.ApprovalTemplateDetail{}).
		Where("approval_template_header_id = ?", approvalHeader.ApprovalTemplateHeaderId).
		Where("approver_by = ?", approvalDetail.ApproverBy).
		First(&approvalTemplateDetailData).Error; err != nil {
		return err
	}

	existedApprovalDetail := models.ApprovalDetail{}
	if err := a.db.Model(&models.ApprovalDetail{}).
		Where("approval_header_id = ?", approvalHeader.ApprovalHeaderId).
		Where("approver_by = ?", approvalDetail.ApproverBy).
		First(&existedApprovalDetail).Error; err != nil {
		return a.db.Model(models.ApprovalDetail{}).Create(&approvalDetail).Error
	}

	return a.db.Model(models.ApprovalDetail{}).
		Where("approval_header_id = ?", approvalHeader.ApprovalHeaderId).
		Where("approver_by = ?", approvalDetail.ApproverBy).
		Updates(map[string]interface{}{
			"approval_status": approvalDetail.ApprovalStatus,
			"approver_by":     approvalDetail.ApproverBy,
			"remark":          approvalDetail.Remark,
		}).Error
}

func CreateApproval(db *gorm.DB, createApprovalDto approval.CreateApprovalDto) (models.ApprovalHeader, error) {
	approvalTemplateHeader := models.ApprovalTemplateHeader{}

	if err := db.Model(models.ApprovalTemplateHeader{}).
		Where("template_type = ?", createApprovalDto.TemplateType).
		First(&approvalTemplateHeader).Error; err != nil {
		return models.ApprovalHeader{}, err
	}

	approval, err := mappers.ToApprovalHeader(createApprovalDto, approvalTemplateHeader)
	if err != nil {
		return models.ApprovalHeader{}, err
	}

	if err := db.Create(&approval).Error; err != nil {
		return models.ApprovalHeader{}, err
	}
	return approval, nil
}

func (a *approvalRepository) Delete(approvalHeaderId string) error {
	return a.db.Delete(&models.ApprovalHeader{}, "approval_header_id = ?", approvalHeaderId).Error
}

func (a *approvalRepository) Detail(approvalHeaderId string) (response.PaginateResponseDto[[]approval.ApprovalDetailResponseDto], error) {
	var data []approval.ApprovalDetailResponseDto
	var totalRecord int64
	var totalPage int

	err := a.db.
		Table("hrms_approval_template_detail atd").
		Select(`
			atd.approval_template_detail_id,
			atd.approval_template_header_id,
			atd.approver_by,
			atd.sequence_number,

			ad.approval_detail_id,
			ad.approval_status,
			ad.remark
		`).
		Joins(`
			LEFT JOIN hrms_approval_detail ad 
			ON ad.approver_by = atd.approver_by 
			AND ad.approval_header_id = ?
		`, approvalHeaderId).
		Where("atd.approval_template_header_id = (?)",
			a.db.
				Table("hrms_approval_header").
				Select("approval_template_header_id").
				Where("approval_header_id = ?", approvalHeaderId),
		).
		Order("atd.sequence_number ASC").
		Scan(&data).Error

	dataAkhir := response.PaginateResponseDto[[]approval.ApprovalDetailResponseDto]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	return dataAkhir, err
}

func (a *approvalRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalHeader], error) {
	var data []models.ApprovalHeader
	var totalRecord int64
	var totalPage int

	modelDb := a.db.Model(&models.ApprovalHeader{})

	dataAkhir := response.PaginateResponseDto[[]models.ApprovalHeader]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "approval_header_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			approval_header_id::text LIKE ? OR
			approval_template_header_id::text LIKE ? OR
			approval_doc_id::text LIKE ? OR
			requester_by::text LIKE ? OR
			created_by LIKE ? OR
			updated_by LIKE ?
		`,
			search, search, search,
			search, search, search,
		)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"approval_header_id": {
			Field: "approval_header_id",
			Query: " = ?",
		},
		"approval_template_header_id": {
			Field: "approval_template_header_id",
			Query: " = ?",
		},
		"approval_doc_id": {
			Field: "approval_doc_id",
			Query: " = ?",
		},
		"requester_by": {
			Field: "requester_by",
			Query: " = ?",
		},
		"created_by": {
			Field: "created_by",
			Query: " LIKE ?",
		},
		"updated_by": {
			Field: "updated_by",
			Query: " LIKE ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Joins("ApprovalTemplateHeader").Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db}
}
