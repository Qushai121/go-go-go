package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApprovalRepository interface {
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalHeader], error)
	Detail(approvalHeaderId string) (response.PaginateResponseDto[models.ApprovalHeader], error)
	Delete(approvalHeaderId string) error
	Approve(approvalDetail models.ApprovalDetail) error
}

type approvalRepository struct {
	db *gorm.DB
}

func (a *approvalRepository) Approve(approvalDetail models.ApprovalDetail) error {
	return a.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "approval_header_id"},
			{Name: "approver_by"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"approval_status": approvalDetail.ApprovalStatus,
			"updated_by":      approvalDetail.UpdatedBy,
			"updated_at":      gorm.Expr("NOW()"),
		}),
	}).Create(&approvalDetail).Error

}

func (a *approvalRepository) Delete(approvalHeaderId string) error {
	return a.db.Delete(&models.ApprovalHeader{}, "approval_header_id = ?", approvalHeaderId).Error
}

func (a *approvalRepository) Detail(approvalHeaderId string) (response.PaginateResponseDto[models.ApprovalHeader], error) {
	var data models.ApprovalHeader
	var totalRecord int64
	var totalPage int

	modelDb := a.db.Model(&models.ApprovalHeader{})

	dataAkhir := response.PaginateResponseDto[models.ApprovalHeader]{
		Data:        	data,
		TotalRecord: 	totalRecord,
		TotalPage: 		totalPage,
	}

	err := modelDb.Joins("ApprovalTemplateHeader").Find(&data).Error
	return dataAkhir,err
}

func (a *approvalRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ApprovalHeader], error) {
	var data []models.ApprovalHeader
	var totalRecord int64
	var totalPage int

	modelDb := a.db.Model(&models.ApprovalHeader{})

	dataAkhir := response.PaginateResponseDto[[]models.ApprovalHeader]{
		Data:        	data,
		TotalRecord: 	totalRecord,
		TotalPage: 		totalPage,
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

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Joins("ApprovalTemplateHeader").Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db}
}
