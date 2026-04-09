package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(shift *models.Shift) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Shift], error)
	Delete(shiftId string) error
	Update(shift *models.Shift) error
}

type shiftRepository struct {
	db *gorm.DB
}

func (s *shiftRepository) Create(shift *models.Shift) error {
	return s.db.Create(shift).Error
}

func (s *shiftRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Shift], error) {
	var data = []models.Shift{}
	modelDb := s.db.Model(&models.Shift{})
	var totalRecord int64
	var totalPage int

	dataAkhir := response.PaginateResponseDto[[]models.Shift]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "shift_id"
		queryParams.SortBy = &sort
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (c *shiftRepository) Delete(shiftId string) error {
	return c.db.Delete(&models.Shift{}, "shift_id = ?", shiftId).Error
}

func (c *shiftRepository) Update(shift *models.Shift) error {
	return c.db.Model(&models.Shift{}).Where("shift_id = ?",shift.ShiftId).Updates(&shift).Error
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db}
}
