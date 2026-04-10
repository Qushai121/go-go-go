package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type WFHRepository interface {
	Create(data *models.WFH) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.WFH], error)
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.WFH], error)
	Delete(id string) error
	Update(data *models.WFH) error
}

type wfhRepository struct {
	db *gorm.DB
}


func NewWFHRepository(db *gorm.DB) WFHRepository {
	return &wfhRepository{db}
}


func (r *wfhRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.WFH], error) {
	var data []models.WFH
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.WFH{})

	// default sort
	if queryParams.SortBy == nil {
		sort := "wfh_id"
		queryParams.SortBy = &sort
	}

	result := response.PaginateResponseDto[[]models.WFH]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	err := utils.GetQuery(queryParams, modelDb, &result.TotalRecord, &result.TotalPage).
		Where("user_id = ?", userId).
		Find(&result.Data).Error

	return result, err
}

func (r *wfhRepository) Create(data *models.WFH) error {
	return r.db.Create(data).Error
}

func (r *wfhRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.WFH], error) {
	var data []models.WFH
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.WFH{})

	// default sort
	if queryParams.SortBy == nil {
		sort := "wfh_id"
		queryParams.SortBy = &sort
	}

	result := response.PaginateResponseDto[[]models.WFH]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.Search != nil{
		search := "%" + *queryParams.Search + "%"
		modelDb = modelDb.Where(`
				wfh_id::text LIKE ? OR
				user_id::text LIKE ? OR
				remarks LIKE ? OR
				start_time::text LIKE ? OR
				end_time::text LIKE ?
			`, search, search, search, search, search)
	}
	
	err := utils.GetQuery(queryParams, modelDb, &result.TotalRecord, &result.TotalPage).Find(&result.Data).Error

	return result, err
}

func (r *wfhRepository) Update(data *models.WFH) error {
	return r.db.Model(&models.WFH{}).
		Where("wfh_id = ?", data.WFHId).
		Updates(data).Error
}

func (r *wfhRepository) Delete(id string) error {
	return r.db.Delete(&models.WFH{}, "wfh_id = ?", id).Error
}
