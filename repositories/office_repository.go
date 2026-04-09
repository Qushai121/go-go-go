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
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "office_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func NewOfficeRepository(db *gorm.DB) OfficeRepository {
	return &officeRepository{db}
}
