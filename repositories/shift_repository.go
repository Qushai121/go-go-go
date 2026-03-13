package repositories

import (
	"hrms_go/models"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(shift *models.Shift) error
	FindAll() ([]models.Shift, error)
}

type shiftRepository struct {
	db *gorm.DB
}

func (s *shiftRepository) Create(shift *models.Shift) error {
	return s.db.Create(shift).Error;
}

func (s *shiftRepository) FindAll() ([]models.Shift, error) {
	data := []models.Shift{}
	err := s.db.Find(&data).Error
	return data, err
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db}
}
