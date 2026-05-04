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
		TotalPage:   totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "shift_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			shift_id::text LIKE ? OR
			shift_code LIKE ? OR
			shift_name LIKE ? OR
			shift_duration::text LIKE ? OR
			start_time LIKE ? OR
			end_time LIKE ? OR
			grace_period::text LIKE ?
		`, search, search, search, search, search, search, search)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"shift_id": {
			Field: "shift_id",
			Query: " = ?",
		},
		"shift_code": {
			Field: "shift_code",
			Query: " LIKE ?",
		},
		"shift_name": {
			Field: "shift_name",
			Query: " LIKE ?",
		},
		"shift_duration": {
			Field: "shift_duration",
			Query: " = ?",
		},
		"start_time": {
			Field: "start_time",
			Query: " >= ?",
		},
		"end_time": {
			Field: "end_time",
			Query: " <= ?",
		},
		"grace_period": {
			Field: "grace_period",
			Query: " = ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &dataAkhir.TotalRecord, &dataAkhir.TotalPage, &allowedDynamicList).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (c *shiftRepository) Delete(shiftId string) error {
	return c.db.Delete(&models.Shift{}, "shift_id = ?", shiftId).Error
}

func (c *shiftRepository) Update(shift *models.Shift) error {
	return c.db.Model(&models.Shift{}).Where("shift_id = ?", shift.ShiftId).Updates(&shift).Error
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db}
}
