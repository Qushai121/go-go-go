package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Submission], error)
	Create(submission *models.Submission) error
	Update(submission *models.Submission) error
	Delete(id string) error
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Submission], error)
}

type submissionRepository struct {
	db *gorm.DB
}

func (r *submissionRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Submission], error) {
	var data []models.Submission
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Submission{}).
		Where("user_id = ?", userId).
		Joins("Claims")

	dataAkhir := response.PaginateResponseDto[[]models.Submission]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "date"
		queryParams.SortBy = &sort
	}

	allowedDynamicList := submissionDynamicSearchFields()

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (r *submissionRepository) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) Update(submission *models.Submission) error {
	return r.db.Model(&models.Submission{}).Where("submission_id = ?", submission.SubmissionID).Updates(submission).Error
}

func (r *submissionRepository) Delete(id string) error {
	return r.db.Delete(&models.Submission{}, "submission_id = ?", id).Error
}

func (r *submissionRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Submission], error) {
	var data []models.Submission
	var total int64
	var totalPage int

	modelDb := r.db.Model(&models.Submission{}).Joins("Claims")

	dataAkhir := response.PaginateResponseDto[[]models.Submission]{
		Data:        data,
		TotalRecord: total,
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "submission_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			submission_id::text LIKE ? OR
			user_id::text LIKE ? OR
			request_number LIKE ? OR
			submit_date::text LIKE ? OR
			status LIKE ? OR
			remarks LIKE ? OR
			amount::text LIKE ?
		`,
			search, // submission_id
			search, // user_id
			search, // request_number
			search, // submit_date
			search, // status
			search, // remarks
			search, // amount
		)
	}

	allowedDynamicList := submissionDynamicSearchFields()

	if err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error; err != nil {
		return dataAkhir, err
	}

	return dataAkhir, nil
}

func submissionDynamicSearchFields() map[string]dto.DynamicSearchDto {
	return map[string]dto.DynamicSearchDto{
		"submission_id": {
			Field: "hrms_submissions.submission_id",
			Query: " = ?",
		},
		"user_id": {
			Field: "hrms_submissions.user_id",
			Query: " = ?",
		},
		"request_number": {
			Field: "hrms_submissions.request_number",
			Query: " LIKE ?",
		},
		"submit_date": {
			Field: "hrms_submissions.submit_date",
			Query: " >= ?",
		},
		"status": {
			Field: "hrms_submissions.status",
			Query: " LIKE ?",
		},
		"remarks": {
			Field: "hrms_submissions.remarks",
			Query: " LIKE ?",
		},
		"amount": {
			Field: "hrms_submissions.amount",
			Query: " = ?",
		},
	}
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}
