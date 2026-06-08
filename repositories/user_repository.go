package repositories

import (
	"hrms_go/dto"
	"hrms_go/dto/response"
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.User], error)
	FindByID(userId string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Count() (int64, error)
	Update(body *models.User) error
	UpdateProfilePicture(body *models.User) error
	Me(userId string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// Me implements [UserRepository].
func (r *userRepository) Me(userId string) (*models.User, error) {
	var data = models.User{}
	err := r.db.Model(models.User{}).Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return &data, err
	}

	err = r.loadUserMappings(&data)
	return &data, err
}

// Count implements [UserRepository].
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// FindAll implements [UserRepository].
func (r *userRepository) FindAll(queryParams *dto.PaginateFieldDto) (response.PaginateResponseDto[[]models.User], error) {
	var data []models.User
	var totalRecord int64
	var totalPage int

	modelDb := r.db.Model(&models.User{})

	// default sort
	if queryParams.SortBy == nil {
		sort := "user_id"
		queryParams.SortBy = &sort
	}

	result := response.PaginateResponseDto[[]models.User]{
		Data:        data,
		TotalRecord: totalRecord,
		TotalPage:   totalPage,
	}

	if queryParams.Search != nil && *queryParams.Search != "" && (queryParams.DynamicFieldSearch == nil || *queryParams.DynamicFieldSearch == "") {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			user_id::text LIKE ? OR
			employee_nik LIKE ? OR
			fullname LIKE ? OR
			email LIKE ? OR
			company_code LIKE ? OR
			division_code LIKE ? OR
			department_code LIKE ? OR
			title_code LIKE ? OR
			is_active LIKE ? OR
			is_locked LIKE ? OR
			created_at::text LIKE ? OR
			updated_at::text LIKE ?
		`,
			search, search, search, search, search, search,
			search, search, search, search, search, search,
		)
	}

	allowedDynamicList := map[string]dto.DynamicSearchDto{
		"user_id": {
			Field: "user_id",
			Query: " = ?",
		},
		"employee_nik": {
			Field: "employee_nik",
			Query: " LIKE ?",
		},
		"fullname": {
			Field: "fullname",
			Query: " LIKE ?",
		},
		"email": {
			Field: "email",
			Query: " LIKE ?",
		},
		"company_code": {
			Field: "company_code",
			Query: " LIKE ?",
		},
		"division_code": {
			Field: "division_code",
			Query: " LIKE ?",
		},
		"department_code": {
			Field: "department_code",
			Query: " LIKE ?",
		},
		"title_code": {
			Field: "title_code",
			Query: " LIKE ?",
		},
		"is_active": {
			Field: "is_active",
			Query: " = ?",
		},
		"is_locked": {
			Field: "is_locked",
			Query: " = ?",
		},
		"need_reset": {
			Field: "need_reset",
			Query: " = ?",
		},
	}

	err := utils.GetQueryBase(queryParams, modelDb, &result.TotalRecord, &result.TotalPage, &allowedDynamicList).Find(&result.Data).Error

	return result, err
}

func (r *userRepository) FindByID(userId string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_id = ?", userId).First(&user).Error
	return &user, err
}

// FindByEmail implements [UserRepository].
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return &user, err
	}

	err = r.loadUserMappings(&user)
	return &user, err
}

func (r *userRepository) loadUserMappings(user *models.User) error {
	if err := r.db.Table("hrms_user_company uc").
		Select(`
			uc.*,
			c.company_id AS company_id,
			c.company_name AS company_name
		`).
		Joins("LEFT JOIN hrms_company c ON c.company_code = uc.company_code").
		Where("uc.employee_nik = ?", user.EmployeeNIK).
		Find(&user.UserCompany).Error; err != nil {
		return err
	}

	if err := r.db.Table("hrms_user_site us").
		Select(`
			us.*,
			s.site_id AS site_id,
			s.site_name AS site_name,
			s.site_phone AS site_phone,
			s.site_address AS site_address,
			s.site_latitude AS site_latitude,
			s.site_longitude AS site_longitude,
			s.max_radius AS max_radius
		`).
		Joins(`
			LEFT JOIN hrms_site s
				ON s.company_code = us.company_code
				AND s.site_type = us.site_type
				AND s.site_code = us.site_code
		`).
		Where("us.employee_nik = ?", user.EmployeeNIK).
		Order("us.site_type ASC, s.site_name ASC").
		Find(&user.UserSite).Error; err != nil {
		return err
	}

	if err := r.db.Table("hrms_employee_shift_period hsep").
		Select(`
		hsep.week_start_date AS week_start_date,
		hsep.week_end_date AS week_end_date,

		(hsep.week_start_date + (hes.weekday_id - 1))::date AS shift_date,

		hes.employee_nik AS employee_nik,
		hs.shift_name AS shift_name,
		hs.is_active AS is_active,
		hs.working_hours_type AS working_hours_type,
		hes.weekday_id AS weekday_id,
		hes.weekday_name AS weekday_name,
		hes.event_code AS event_code,
		hes.event_name AS event_name,
		LEFT(hes.start_time, 5) AS start_time,
		LEFT(hes.end_time, 5) AS end_time,
		hes.tolerance_before AS tolerance_before,
		hes.tolarance_after AS tolarance_after
	`).
		Joins(`
		LEFT JOIN hrms_employee_shift hes
			ON hsep.weekly_id = hes.weekly_id
	`).
		Joins(`
		LEFT JOIN hrms_shift hs
			ON hes.shift_id = hs.shift_id
	`).
		Where("hes.employee_nik = ?", user.EmployeeNIK).
		Where("CURRENT_DATE BETWEEN hsep.week_start_date AND hsep.week_end_date").
		Order(`
		shift_date ASC,
		hes.start_time ASC,
		hes.event_code ASC
	`).
		Find(&user.UserShift).Error; err != nil {
		return err
	}

	if err := r.db.Table("hrms_leave_balance hlb").
		Select(`
		hlb.*
	`).
		Where("hlb.employee_nik = ?", user.EmployeeNIK).
		Where("CURRENT_DATE BETWEEN hlb.leave_period_start AND hlb.leave_period_end").
		Order("hlb.leave_period_start DESC").
		Find(&user.UserLeaveBalance).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Update(body *models.User) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", body.UserId).Updates(body).Error
}

func (r *userRepository) UpdateProfilePicture(body *models.User) error {
	var user models.User

	if err := r.db.
		Select("profile_picture_url").
		Where("user_id = ?", body.UserId).
		First(&user).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.User{}).
		Where("user_id = ?", body.UserId).
		Updates(map[string]interface{}{
			"profile_picture_url": body.ProfilePictureUrl,
			"updated_by":          body.UpdatedBy,
		}).Error; err != nil {
		return err
	}

	if user.ProfilePictureUrl != "" && user.ProfilePictureUrl != body.ProfilePictureUrl {
		_ = utils.RemoveFileFromPath(user.ProfilePictureUrl)
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
