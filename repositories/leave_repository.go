package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type LeaveRepository interface {
	Create(leave *models.Leave) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Leave], error)
	Delete(leaveId string) error
	Update(leave *models.Leave) error
}

type leaveRepository struct {
	db *gorm.DB
}

func (r *leaveRepository) Create(leave *models.Leave) error {
	return r.db.Create(leave).Error
}

func (r *leaveRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.Leave], error) {
	var data []models.Leave
	modelDb := r.db.Model(&models.Leave{})
	var totalRecord int64
	var totalPage int


	dataAkhir := response.PaginateResponseDto[[]models.Leave]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage: totalPage,
	}

	if queryParams.SortBy == nil {
		sort := "leave_id"
		queryParams.SortBy = &sort
	}

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			leave_id::text LIKE ? OR
			request_number LIKE ? OR
			employee_name LIKE ? OR
			leave_type LIKE ? OR
			start_date::text LIKE ? OR
			end_date::text LIKE ? OR
			req_date::text LIKE ? OR
			status LIKE ? OR
			message LIKE ? OR
			cancellation_reason LIKE ? OR
			leave_balance::text LIKE ?
		`,
			search, search, search, search, search,
			search, search, search, search, search,
			search,
		)
	}

	err := utils.GetQuery(queryParams, modelDb, &dataAkhir.TotalRecord,&dataAkhir.TotalPage).Find(&dataAkhir.Data).Error
	return dataAkhir, err
}

func (r *leaveRepository) Update(leave *models.Leave) error {
	return r.db.Model(&models.Leave{}).Where("leave_id = ?", leave.LeaveID).Updates(&leave).Error
}

func (r *leaveRepository) Delete(leaveId string) error {
	return r.db.Delete(&models.Leave{}, "leave_id = ?", leaveId).Error
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{db}
}
