package repositories

import (
	"strings"

	"hrms_go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeaveTypeRepository interface {
	FindAll(filter models.LeaveTypeFilter) ([]models.LeaveType, error)
	FindByID(id uuid.UUID) (*models.LeaveType, error)
}

type leaveTypeRepository struct {
	db *gorm.DB
}

func NewLeaveTypeRepository(db *gorm.DB) LeaveTypeRepository {
	return &leaveTypeRepository{db: db}
}

func (r *leaveTypeRepository) FindAll(filter models.LeaveTypeFilter) ([]models.LeaveType, error) {
	var data []models.LeaveType

	query := r.db.Model(&models.LeaveType{})

	companyCode := strings.TrimSpace(filter.CompanyCode)
	search := strings.TrimSpace(filter.Search)
	isActive := strings.TrimSpace(filter.IsActive)
	deductLeaveBalance := strings.TrimSpace(filter.DeductLeaveBalance)

	if companyCode != "" {
		query = query.Where("company_code = ?", companyCode)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(`
			leave_type_code ILIKE ?
			OR leave_type_name ILIKE ?
			OR COALESCE(leave_type_description, '') ILIKE ?
		`, like, like, like)
	}

	if isActive != "" {
		query = query.Where("COALESCE(is_active, true) = ?", parseBoolFilter(isActive))
	} else {
		query = query.Where("COALESCE(is_active, true) = true")
	}

	if deductLeaveBalance != "" {
		query = query.Where("COALESCE(deduct_leave_balance, false) = ?", parseBoolFilter(deductLeaveBalance))
	}

	err := query.
		Order("sort_order ASC, leave_type_name ASC").
		Find(&data).Error

	return data, err
}

func (r *leaveTypeRepository) FindByID(id uuid.UUID) (*models.LeaveType, error) {
	var data models.LeaveType
	err := r.db.Where("leave_type_id = ?", id).First(&data).Error
	return &data, err
}

func parseBoolFilter(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1" || value == "y" || value == "yes"
}
