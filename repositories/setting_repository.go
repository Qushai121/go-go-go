package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type SettingRepository interface {
	Create(setting *models.Setting) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Setting], error)
	Update(setting *models.Setting) error
	Delete(settingId string) error
}

type settingRepository struct {
	db *gorm.DB
}

func (r *settingRepository) Create(setting *models.Setting) error {
	return r.db.Create(setting).Error
}

func (r *settingRepository) Update(setting *models.Setting) error {
	return r.db.Model(&models.Setting{}).
		Where("setting_id = ?", setting.SettingId).
		Updates(setting).Error
}

func (r *settingRepository) Delete(settingId string) error {
	return r.db.Delete(&models.Setting{}, "setting_id = ?", settingId).Error
}

func (r *settingRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Setting], error) {
	var data []models.Setting
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.Setting{})

	dataAkhir := response.PaginateResponseDto[[]models.Setting]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "setting_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error

	return dataAkhir, err
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db}
}
