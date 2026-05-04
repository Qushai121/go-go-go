package repositories

import (
	"hrms_go/constant"
	"hrms_go/dto"
	"hrms_go/dto/approval"
	"hrms_go/dto/response"
	"hrms_go/dto/wfh"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type WFHRepository interface {
	Create(data *models.WFH) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]wfh.WFHApproval], error)
	FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]wfh.WFHApproval], error)
	Delete(id string) error
	Update(data *models.WFH) error
}

type wfhRepository struct {
	db *gorm.DB
}

func (r *wfhRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]wfh.WFHApproval], error) {
	return r.findAll(nil, queryParams)
}

func (r *wfhRepository) FindByUser(userId string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]wfh.WFHApproval], error) {
	return r.findAll(&userId, queryParams)
}

func (r *wfhRepository) findAll(userId *string, queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]wfh.WFHApproval], error) {
	var data []wfh.WFHApproval
	var totalRecord int64
	var totalPage int

	result := response.PaginateResponseDto[[]wfh.WFHApproval]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	modelDb := r.db.
		Table("hrms_wfh w").
		Select(`
			w.wfh_id,
			w.user_id,
			w.remarks,
			ah.approval_header_id,

			CASE
			WHEN COUNT(ad.approval_detail_id) = 0 THEN 'PENDING'

			WHEN COUNT(*) = COUNT(CASE WHEN ad.approval_status = 'approved' THEN 1 END)
				THEN 'APPROVED'

			WHEN COUNT(*) = COUNT(CASE WHEN ad.approval_status = 'pending' THEN 1 END)
				THEN 'PENDING'

			ELSE 'NEED_APPROVAL'
			END as final_status
		`).
		Joins(`
			LEFT JOIN hrms_approval_header ah
			ON ah.approval_doc_id = w.wfh_id
		`).
		Joins(`
			LEFT JOIN hrms_approval_detail ad
			ON ad.approval_header_id = ah.approval_header_id
		`).
		Group(`
			w.wfh_id,
			w.user_id,
			w.remarks,
			ah.approval_header_id
		`)

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			w.wfh_id::text LIKE ? OR
			w.user_id::text LIKE ? OR
			w.remarks LIKE ? OR
			w.start_time::text LIKE ? OR
			w.end_time::text LIKE ? OR
			w.created_at::text LIKE ? OR
			w.updated_at::text LIKE ?
		`,
			search, // wfh_id
			search, // user_id
			search, // remarks
			search, // start_time
			search, // end_time
			search, // created_at
			search, // updated_at
		)
	}

	if userId != nil {
		modelDb = modelDb.Where("w.user_id = ?", *userId)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"wfh_id": {
			Field: "w.wfh_id",
			Query: " = ?",
		},
		"start_time": {
			Field: "w.start_time",
			Query: " > ?",
		},
		"end_time": {
			Field: "w.end_time",
			Query: " < ?",
		},
		"remarks": {
			Field: "w.remarks",
			Query: " LIKE ?",
		},
		"user_id": {
			Field: "w.user_id",
			Query: " = ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error
	return result, err
}

func NewWFHRepository(db *gorm.DB) WFHRepository {
	return &wfhRepository{db}
}

func (r *wfhRepository) Create(data *models.WFH) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		_, err := CreateApproval(r.db, approval.CreateApprovalDto{
			TemplateType:  constant.WFH,
			RequesterBy:   data.CreatedBy,
			ApprovalDocId: data.WFHId,
		})

		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// func (r *wfhRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.WFH], error) {
// 	var data []models.WFH
// 	var totalRecord int64
// 	var totalPage int

// 	modelDb := r.db.Model(&models.WFH{})

// 	// default sort
// 	if queryParams.SortBy == nil {
// 		sort := "wfh_id"
// 		queryParams.SortBy = &sort
// 	}

// 	result := response.PaginateResponseDto[[]models.WFH]{
// 		Data:        data,
// 		TotalRecord: totalRecord,
// 		TotalPage:   totalPage,
// 	}

// 	if queryParams.Search != nil{
// 		search := "%" + *queryParams.Search + "%"
// 		modelDb = modelDb.Where(`
// 				wfh_id::text LIKE ? OR
// 				user_id::text LIKE ? OR
// 				remarks LIKE ? OR
// 				start_time::text LIKE ? OR
// 				end_time::text LIKE ?
// 			`, search, search, search, search, search)
// 	}

// 	err := utils.GetQuery(queryParams, modelDb, &result.TotalRecord, &result.TotalPage).Find(&result.Data).Error

// 	return result, err
// }

func (r *wfhRepository) Update(data *models.WFH) error {
	return r.db.Model(&models.WFH{}).
		Where("wfh_id = ?", data.WFHId).
		Updates(data).Error
}

func (r *wfhRepository) Delete(id string) error {
	return r.db.Delete(&models.WFH{}, "wfh_id = ?", id).Error
}
