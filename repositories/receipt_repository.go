package repositories

import (
	"errors"
	"fmt"
	"hrms_go/constant"
	"hrms_go/dto"
	"hrms_go/dto/approval"
	receiptDto "hrms_go/dto/receipt"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/models/base"
	"hrms_go/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReceiptRepository interface {
	Create(req receiptDto.CreateReceiptRequestDto, createdBy string, requesterBy string) (models.ReceiptHeader, error)
	Submit(req receiptDto.SubmitReceiptRequestDto, createdBy string, requesterBy string) (models.ReceiptHeader, error)
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ReceiptHeader], error)
	FindDraftDetails(queryParams *dto.PaginateFieldDto, createdBy string) (response.PaginateResponseDto[[]models.ReceiptDetail], error)
	Detail(receiptId string) (models.ReceiptHeader, error)
	UpdateHeader(receiptId string, req receiptDto.UpdateReceiptHeaderDto, updatedBy string) (models.ReceiptHeader, error)
	UpdateSubmission(receiptId string, req receiptDto.UpdateReceiptSubmissionDto, updatedBy string) (models.Submission, error)
	Delete(receiptId string) error
	CreateDraftDetail(req receiptDto.CreateReceiptDetailDto, createdBy string) (models.ReceiptDetail, error)
	CreateDetail(receiptId string, req receiptDto.CreateReceiptDetailDto, createdBy string) (models.ReceiptDetail, error)
	UpdateDetail(receiptDetailId string, req receiptDto.UpdateReceiptDetailDto, updatedBy string) (models.ReceiptDetail, error)
	DeleteDetail(receiptDetailId string) error
}

type receiptRepository struct {
	db *gorm.DB
}

func (r *receiptRepository) Create(req receiptDto.CreateReceiptRequestDto, createdBy string, requesterBy string) (models.ReceiptHeader, error) {
	var header models.ReceiptHeader

	err := r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		receiptDate := now
		if req.ReceiptCreateDate != nil {
			receiptDate = *req.ReceiptCreateDate
		}

		totalAmount := 0.0
		for _, detail := range req.Details {
			totalAmount += detail.ReceiptAmount
		}

		header = models.ReceiptHeader{
			ReceiptId:          uuid.New(),
			EmployeeNik:        req.EmployeeNik,
			ReceiptCreateDate:  receiptDate,
			TotalReceipt:       len(req.Details),
			TotalAmountReceipt: totalAmount,
			ObjectCode:         "RECEIPT_HEADER",
			AuditFields: base.AuditFields{
				CreatedBy: createdBy,
			},
		}

		if err := tx.Create(&header).Error; err != nil {
			return err
		}

		details := make([]models.ReceiptDetail, 0, len(req.Details))
		for _, item := range req.Details {
			details = append(details, models.ReceiptDetail{
				ReceiptDetailId:    uuid.New(),
				ReceiptId:          header.ReceiptId,
				ReceiptDate:        item.ReceiptDate,
				ReceiptType:        item.ReceiptType,
				ReceiptAmount:      item.ReceiptAmount,
				ReceiptDescription: item.ReceiptDescription,
				ReceiptImage:       item.ReceiptImage,
				ObjectCode:         "RECEIPT_DETAIL",
				AuditFields: base.AuditFields{
					CreatedBy: createdBy,
				},
			})
		}

		if err := tx.Create(&details).Error; err != nil {
			return err
		}

		submissionNumber, err := r.generateSubmissionNumber(tx, now)
		if err != nil {
			return err
		}

		submission := models.Submission{
			SubmissionId:     uuid.New(),
			ReceiptId:        header.ReceiptId,
			SubmissionNumber: submissionNumber,
			SubmissionDate:   now,
			Status:           "P",
			CurrentStep:      1,
			ObjectCode:       "SUBMISSION",
			AuditFields: base.AuditFields{
				CreatedBy: createdBy,
			},
		}
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}

		if requesterBy != "" {
			approvalHeader, err := CreateApproval(tx, approval.CreateApprovalDto{
				TemplateType:  constant.RECEIPT,
				RequesterBy:   requesterBy,
				CreatedBy:     createdBy,
				ApprovalDocId: submission.SubmissionId,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err == nil {
				submission.ApprovalHeaderId = &approvalHeader.ApprovalHeaderId
				if err := tx.Model(&models.Submission{}).
					Where("submission_id = ?", submission.SubmissionId).
					Update("approvalheader_id", approvalHeader.ApprovalHeaderId).Error; err != nil {
					return err
				}
			}
		}

		header.ReceiptDetails = details
		header.Submission = &submission
		return nil
	})

	return header, err
}

func (r *receiptRepository) Submit(req receiptDto.SubmitReceiptRequestDto, createdBy string, requesterBy string) (models.ReceiptHeader, error) {
	var header models.ReceiptHeader

	err := r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		receiptDate := now
		if req.ReceiptCreateDate != nil {
			receiptDate = *req.ReceiptCreateDate
		}

		detailIds := make([]uuid.UUID, 0, len(req.ReceiptDetailIds))
		for _, rawId := range req.ReceiptDetailIds {
			parsedId, err := uuid.Parse(rawId)
			if err != nil {
				return err
			}
			detailIds = append(detailIds, parsedId)
		}

		detailQuery := tx.Model(&models.ReceiptDetail{}).
			Joins("LEFT JOIN hrms_receipt_header rh ON rh.receipt_id = hrms_receipt_detail.receipt_id").
			Where("hrms_receipt_detail.receipt_detail_id IN ?", detailIds).
			Where("rh.receipt_id IS NULL")
		if createdBy != "" {
			detailQuery = detailQuery.Where("hrms_receipt_detail.created_by = ?", createdBy)
		}

		var details []models.ReceiptDetail
		if err := detailQuery.Find(&details).Error; err != nil {
			return err
		}

		if len(details) != len(detailIds) {
			return fmt.Errorf("some receipt details are invalid or already submitted")
		}

		totalAmount := 0.0
		for _, detail := range details {
			totalAmount += detail.ReceiptAmount
		}

		header = models.ReceiptHeader{
			ReceiptId:          uuid.New(),
			EmployeeNik:        req.EmployeeNik,
			ReceiptCreateDate:  receiptDate,
			TotalReceipt:       len(details),
			TotalAmountReceipt: totalAmount,
			ObjectCode:         "RECEIPT_HEADER",
			AuditFields: base.AuditFields{
				CreatedBy: createdBy,
			},
		}

		if err := tx.Create(&header).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ReceiptDetail{}).
			Where("receipt_detail_id IN ?", detailIds).
			Updates(map[string]interface{}{
				"receipt_id": header.ReceiptId,
				"updated_by": createdBy,
			}).Error; err != nil {
			return err
		}

		submissionNumber, err := r.generateSubmissionNumber(tx, now)
		if err != nil {
			return err
		}

		submission := models.Submission{
			SubmissionId:     uuid.New(),
			ReceiptId:        header.ReceiptId,
			SubmissionNumber: submissionNumber,
			SubmissionDate:   now,
			Status:           "P",
			CurrentStep:      1,
			ObjectCode:       "SUBMISSION",
			AuditFields: base.AuditFields{
				CreatedBy: createdBy,
			},
		}
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}

		if requesterBy != "" {
			approvalHeader, err := CreateApproval(tx, approval.CreateApprovalDto{
				TemplateType:  constant.RECEIPT,
				RequesterBy:   requesterBy,
				CreatedBy:     createdBy,
				ApprovalDocId: submission.SubmissionId,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err == nil {
				submission.ApprovalHeaderId = &approvalHeader.ApprovalHeaderId
				if err := tx.Model(&models.Submission{}).
					Where("submission_id = ?", submission.SubmissionId).
					Update("approvalheader_id", approvalHeader.ApprovalHeaderId).Error; err != nil {
					return err
				}
			}
		}

		header.ReceiptDetails = details
		for index := range header.ReceiptDetails {
			header.ReceiptDetails[index].ReceiptId = header.ReceiptId
		}
		header.Submission = &submission
		return nil
	})

	return header, err
}

func (r *receiptRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ReceiptHeader], error) {
	var data []models.ReceiptHeader
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.ReceiptHeader{}).
		Select(receiptHeaderSelectColumns()).
		Joins("LEFT JOIN hrms_submission s ON s.receipt_id = hrms_receipt_header.receipt_id")
	dataAkhir := response.PaginateResponseDto[[]models.ReceiptHeader]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "receipt_create_date"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			hrms_receipt_header.receipt_id::text LIKE ? OR
			hrms_receipt_header.employee_nik LIKE ? OR
			s.status LIKE ? OR
			hrms_receipt_header.total_receipt::text LIKE ? OR
			hrms_receipt_header.total_amount_receipt::text LIKE ? OR
			s.submission_number LIKE ? OR
			hrms_receipt_header.created_by LIKE ? OR
			hrms_receipt_header.updated_by LIKE ?
		`, search, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"receipt_id":           {Field: "hrms_receipt_header.receipt_id", Query: " = ?"},
		"employee_nik":         {Field: "hrms_receipt_header.employee_nik", Query: " LIKE ?"},
		"receipt_create_date":  {Field: "hrms_receipt_header.receipt_create_date", Query: " >= ?"},
		"total_receipt":        {Field: "hrms_receipt_header.total_receipt", Query: " = ?"},
		"total_amount_receipt": {Field: "hrms_receipt_header.total_amount_receipt", Query: " = ?"},
		"status":               {Field: "s.status", Query: " = ?"},
		"current_step":         {Field: "s.current_step", Query: " = ?"},
		"approvalheader_id":    {Field: "s.approvalheader_id", Query: " = ?"},
		"submission_number":    {Field: "s.submission_number", Query: " LIKE ?"},
		"created_by":           {Field: "hrms_receipt_header.created_by", Query: " LIKE ?"},
		"updated_by":           {Field: "hrms_receipt_header.updated_by", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Preload("Submission").
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func (r *receiptRepository) FindDraftDetails(queryParams *dto.PaginateFieldDto, createdBy string) (response.PaginateResponseDto[[]models.ReceiptDetail], error) {
	var data []models.ReceiptDetail
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.ReceiptDetail{}).
		Select(receiptDetailSelectColumns()).
		Joins("LEFT JOIN hrms_receipt_header rh ON rh.receipt_id = hrms_receipt_detail.receipt_id").
		Where("rh.receipt_id IS NULL")

	if createdBy != "" {
		modelDb = modelDb.Where("hrms_receipt_detail.created_by = ?", createdBy)
	}

	dataAkhir := response.PaginateResponseDto[[]models.ReceiptDetail]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "receipt_date"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
			hrms_receipt_detail.receipt_detail_id::text LIKE ? OR
			hrms_receipt_detail.receipt_type LIKE ? OR
			hrms_receipt_detail.receipt_description LIKE ? OR
			hrms_receipt_detail.receipt_amount::text LIKE ?
		`, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"receipt_detail_id": {Field: "hrms_receipt_detail.receipt_detail_id", Query: " = ?"},
		"receipt_type":      {Field: "hrms_receipt_detail.receipt_type", Query: " LIKE ?"},
		"receipt_date":      {Field: "hrms_receipt_detail.receipt_date", Query: " >= ?"},
		"receipt_amount":    {Field: "hrms_receipt_detail.receipt_amount", Query: " = ?"},
		"created_by":        {Field: "hrms_receipt_detail.created_by", Query: " LIKE ?"},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).
		Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func (r *receiptRepository) Detail(receiptId string) (models.ReceiptHeader, error) {
	var data models.ReceiptHeader
	err := r.db.Preload("ReceiptDetails").Preload("Submission").
		First(&data, "receipt_id = ?", receiptId).Error
	return data, err
}

func (r *receiptRepository) UpdateHeader(receiptId string, req receiptDto.UpdateReceiptHeaderDto, updatedBy string) (models.ReceiptHeader, error) {
	var data models.ReceiptHeader
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&data, "receipt_id = ?", receiptId).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"employee_nik":        req.EmployeeNik,
			"receipt_create_date": req.ReceiptCreateDate,
			"updated_by":          updatedBy,
		}

		if err := tx.Model(&models.ReceiptHeader{}).Where("receipt_id = ?", receiptId).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Preload("ReceiptDetails").Preload("Submission").First(&data, "receipt_id = ?", receiptId).Error
	})
	return data, err
}

func (r *receiptRepository) UpdateSubmission(receiptId string, req receiptDto.UpdateReceiptSubmissionDto, updatedBy string) (models.Submission, error) {
	var data models.Submission
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&data, "receipt_id = ?", receiptId).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":       req.Status,
			"current_step": req.CurrentStep,
			"updated_by":   updatedBy,
		}

		if req.ApprovalHeaderId != nil && *req.ApprovalHeaderId != "" {
			approvalHeaderId, err := uuid.Parse(*req.ApprovalHeaderId)
			if err != nil {
				return err
			}
			updates["approvalheader_id"] = approvalHeaderId
		}

		if err := tx.Model(&models.Submission{}).Where("receipt_id = ?", receiptId).Updates(updates).Error; err != nil {
			return err
		}
		return tx.First(&data, "receipt_id = ?", receiptId).Error
	})
	return data, err
}

func (r *receiptRepository) Delete(receiptId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.ReceiptDetail{}, "receipt_id = ?", receiptId).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Submission{}, "receipt_id = ?", receiptId).Error; err != nil {
			return err
		}
		return tx.Delete(&models.ReceiptHeader{}, "receipt_id = ?", receiptId).Error
	})
}

func (r *receiptRepository) CreateDraftDetail(req receiptDto.CreateReceiptDetailDto, createdBy string) (models.ReceiptDetail, error) {
	detail := models.ReceiptDetail{
		ReceiptDetailId:    uuid.New(),
		ReceiptId:          uuid.New(),
		ReceiptDate:        req.ReceiptDate,
		ReceiptType:        req.ReceiptType,
		ReceiptAmount:      req.ReceiptAmount,
		ReceiptDescription: req.ReceiptDescription,
		ReceiptImage:       req.ReceiptImage,
		ObjectCode:         "RECEIPT_DETAIL",
		AuditFields: base.AuditFields{
			CreatedBy: createdBy,
		},
	}

	return detail, r.db.Create(&detail).Error
}

func (r *receiptRepository) CreateDetail(receiptId string, req receiptDto.CreateReceiptDetailDto, createdBy string) (models.ReceiptDetail, error) {
	detail := models.ReceiptDetail{}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		header := models.ReceiptHeader{}
		if err := tx.First(&header, "receipt_id = ?", receiptId).Error; err != nil {
			return err
		}

		parsedReceiptId, err := uuid.Parse(receiptId)
		if err != nil {
			return err
		}

		detail = models.ReceiptDetail{
			ReceiptDetailId:    uuid.New(),
			ReceiptId:          parsedReceiptId,
			ReceiptDate:        req.ReceiptDate,
			ReceiptType:        req.ReceiptType,
			ReceiptAmount:      req.ReceiptAmount,
			ReceiptDescription: req.ReceiptDescription,
			ReceiptImage:       req.ReceiptImage,
			ObjectCode:         "RECEIPT_DETAIL",
			AuditFields: base.AuditFields{
				CreatedBy: createdBy,
			},
		}
		if err := tx.Create(&detail).Error; err != nil {
			return err
		}
		return r.refreshHeaderTotal(tx, receiptId)
	})
	return detail, err
}

func (r *receiptRepository) UpdateDetail(receiptDetailId string, req receiptDto.UpdateReceiptDetailDto, updatedBy string) (models.ReceiptDetail, error) {
	var detail models.ReceiptDetail
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&detail, "receipt_detail_id = ?", receiptDetailId).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ReceiptDetail{}).Where("receipt_detail_id = ?", receiptDetailId).Updates(map[string]interface{}{
			"receipt_date":        req.ReceiptDate,
			"receipt_type":        req.ReceiptType,
			"receipt_amount":      req.ReceiptAmount,
			"receipt_description": req.ReceiptDescription,
			"receipt_image":       req.ReceiptImage,
			"updated_by":          updatedBy,
		}).Error; err != nil {
			return err
		}

		if err := r.refreshHeaderTotal(tx, detail.ReceiptId.String()); err != nil {
			return err
		}
		return tx.First(&detail, "receipt_detail_id = ?", receiptDetailId).Error
	})
	return detail, err
}

func (r *receiptRepository) DeleteDetail(receiptDetailId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var detail models.ReceiptDetail
		if err := tx.First(&detail, "receipt_detail_id = ?", receiptDetailId).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.ReceiptDetail{}, "receipt_detail_id = ?", receiptDetailId).Error; err != nil {
			return err
		}
		return r.refreshHeaderTotal(tx, detail.ReceiptId.String())
	})
}

func (r *receiptRepository) generateSubmissionNumber(tx *gorm.DB, date time.Time) (string, error) {
	prefix := fmt.Sprintf("SUB-%s-", date.Format("20060102"))
	var total int64
	if err := tx.Model(&models.Submission{}).Where("submission_number LIKE ?", prefix+"%").Count(&total).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%04d", prefix, total+1), nil
}

func (r *receiptRepository) refreshHeaderTotal(tx *gorm.DB, receiptId string) error {
	var totalReceipt int64
	var totalAmount float64
	if err := tx.Model(&models.ReceiptDetail{}).Where("receipt_id = ?", receiptId).Count(&totalReceipt).Error; err != nil {
		return err
	}
	if err := tx.Model(&models.ReceiptDetail{}).Where("receipt_id = ?", receiptId).Select("COALESCE(SUM(receipt_amount), 0)").Scan(&totalAmount).Error; err != nil {
		return err
	}
	return tx.Model(&models.ReceiptHeader{}).Where("receipt_id = ?", receiptId).Updates(map[string]interface{}{
		"total_receipt":        totalReceipt,
		"total_amount_receipt": totalAmount,
	}).Error
}

func receiptHeaderSelectColumns() string {
	return `
		hrms_receipt_header.receipt_id,
		hrms_receipt_header.employee_nik,
		hrms_receipt_header.receipt_create_date,
		hrms_receipt_header.total_receipt,
		hrms_receipt_header.total_amount_receipt,
		hrms_receipt_header.object_code,
		hrms_receipt_header.created_by,
		hrms_receipt_header.created_at,
		hrms_receipt_header.updated_by,
		hrms_receipt_header.updated_at
	`
}

func receiptDetailSelectColumns() string {
	return `
		hrms_receipt_detail.receipt_detail_id,
		hrms_receipt_detail.receipt_id,
		hrms_receipt_detail.receipt_date,
		hrms_receipt_detail.receipt_type,
		hrms_receipt_detail.receipt_amount,
		hrms_receipt_detail.receipt_description,
		hrms_receipt_detail.receipt_image,
		hrms_receipt_detail.object_code,
		hrms_receipt_detail.created_by,
		hrms_receipt_detail.created_at,
		hrms_receipt_detail.updated_by,
		hrms_receipt_detail.updated_at
	`
}

func NewReceiptRepository(db *gorm.DB) ReceiptRepository {
	return &receiptRepository{db}
}
