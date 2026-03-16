package repositories

import (
	"hrms_go/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByEmail(email string) (*models.User, error)
	Count() (int64, error)
	UpdateUserShift(body *models.User) error
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

func (r *userRepository) UpdateUserShift(body *models.User) error  {
	return r.db.Model(&models.User{}).Where("user_id = ?",body.UserId).Updates(body).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
