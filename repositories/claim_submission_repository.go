package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ClaimSubmissionRepository interface {
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ClaimSubmission], error)
	Create(submission *models.ClaimSubmission) error
	Update(submission *models.ClaimSubmission) error
	Delete(id string) error
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ClaimSubmission], error)
}

type claimSubmissionRepository struct {
	db *gorm.DB
}

func (r *claimSubmissionRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ClaimSubmission], error) {
	var data []models.ClaimSubmission
	var totalRecord int64

	modelDb := r.db.Model(&models.ClaimSubmission{}).
		Where("user_id = ?", userId).
		Joins("Claims")

	dataAkhir := response.PaginateResponseDto[[]models.ClaimSubmission]{
		Data:        data,
		TotalRecord: totalRecord,
	}

	if queryParams.SortBy == nil {
		sort := "date"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &totalRecord).Find(&data).Error

	dataAkhir.Data = data
	dataAkhir.TotalRecord = totalRecord

	return dataAkhir, err
}

func NewClaimSubmissionRepository(db *gorm.DB) ClaimSubmissionRepository {
	return &claimSubmissionRepository{db}
}

func (r *claimSubmissionRepository) Create(submission *models.ClaimSubmission) error {
	return r.db.Create(submission).Error
}

func (r *claimSubmissionRepository) Update(submission *models.ClaimSubmission) error {
	return r.db.Model(&models.ClaimSubmission{}).Updates(submission).Error
}

func (r *claimSubmissionRepository) Delete(id string) error {
	return r.db.Delete(&models.ClaimSubmission{}, "submission_id = ?", id).Error
}

func (r *claimSubmissionRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.ClaimSubmission], error) {
	data := []models.ClaimSubmission{}
	modelDb := r.db.Model(&models.ClaimSubmission{}).Joins("Claims")
	var total int64

	dataAkhir := response.PaginateResponseDto[[]models.ClaimSubmission]{
		Data:        data,
		TotalRecord: total,
	}

	if queryParams.SortBy == nil {
		sort := "submission_id"
		queryParams.SortBy = &sort
	}

	if err := utils.GetQuery(queryParams, modelDb, &total).Find(&data).Error; err != nil {
		return dataAkhir, err
	}

	dataAkhir.Data = data
	dataAkhir.TotalRecord = total
	return dataAkhir, nil
}