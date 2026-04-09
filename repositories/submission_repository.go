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
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "date"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (r *submissionRepository) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *submissionRepository) Update(submission *models.Submission) error {
	return r.db.Model(&models.Submission{}).Where("submission_id = ?", submission.SubmissionID).Updates(submission).Error
}

func (r *submissionRepository) Delete(id string) error {
	return r.db.Delete(&models.Submission{}, ).Error
}

func (r *submissionRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Submission], error) {
	data := []models.Submission{}
	modelDb := r.db.Model(&models.Submission{}).Joins("Claims")
	var total int64
	var totalPage int

	dataAkhir := response.PaginateResponseDto[[]models.Submission]{
		Data:        data,
		TotalRecord: total,
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "submission_id"
		queryParams.SortBy = &sort
	}

	if err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error; err != nil {
		return dataAkhir, err
	}

	return dataAkhir, nil
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}
