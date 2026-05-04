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

	if queryParams.Search != nil && *queryParams.Search != "" {
		search := "%" + *queryParams.Search + "%"

		modelDb = modelDb.Where(`
			user_id::text LIKE ? OR
			employee_nik LIKE ? OR
			fullname LIKE ? OR
			email LIKE ? OR
			company_code LIKE ? OR
			branch_code LIKE ? OR
			office_code LIKE ? OR
			division_code LIKE ? OR
			department_code LIKE ? OR
			title_code LIKE ? OR
			is_active LIKE ? OR
			is_locked LIKE ? OR
			created_at::text LIKE ? OR
			updated_at::text LIKE ?
		`,
			search, search, search, search, search, search, search,
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
		"branch_code": {
			Field: "branch_code",
			Query: " LIKE ?",
		},
		"office_code": {
			Field: "office_code",
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

// FindByEmail implements [UserRepository].
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(body *models.User) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", body.UserId).Updates(body).Error
}

func (r *userRepository) UpdateProfilePicture(body *models.User) error {
	var user models.User

	// 1. Get existing user
	if err := r.db.
		Select("profile_picture_url").
		Where("user_id = ?", body.UserId).
		First(&user).Error; err != nil {
		return err
	}

	// 2. Remove old image if exists
	if user.ProfilePictureUrl != "" {
		utils.RemoveFileFromPath(user.ProfilePictureUrl)
	}

	return r.db.Model(&models.User{}).Where("user_id = ?", body.UserId).Updates(body).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
