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
			role LIKE ? OR
			shift_id::text LIKE ? OR
			profile_picture_url LIKE ? OR
			created_at::text LIKE ? OR
			updated_at::text LIKE ?
		`,
			search, // user_id
			search, // employee_nik
			search, // fullname
			search, // email
			search, // role
			search, // shift_id
			search, // profile_picture_url
			search, // created_at
			search, // updated_at
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
		"role": {
			Field: "role",
			Query: " LIKE ?",
		},
		"shift_id": {
			Field: "shift_id",
			Query: " = ?",
		},
		"profile_picture_url": {
			Field: "profile_picture_url",
			Query: " LIKE ?",
		},
		"created_at": {
			Field: "created_at",
			Query: " >= ?",
		},
		"updated_at": {
			Field: "updated_at",
			Query: " <= ?",
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
