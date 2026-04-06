package repositories

import (
	"hrms_go/models"
	"hrms_go/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByEmail(email string) (*models.User, error)
	Count() (int64, error)
	Update(body *models.User) error
	UpdateProfilePicture(body *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

// Count implements [UserRepository].
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// FindAll implements [UserRepository].
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("Shift").Find(&users).Error
	return users, err
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
