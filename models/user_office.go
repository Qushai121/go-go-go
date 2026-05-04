package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type UserOffice struct {
	UserOfficeId uuid.UUID `gorm:"column:user_office_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"user_office_id"`
	CompanyCode  string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	BranchCode   string    `gorm:"column:branch_code;type:varchar(100);not null" json:"branch_code"`
	EmployeeNIK  string    `gorm:"column:employee_nik;type:varchar(50);not null" json:"employee_nik"`
	OfficeCode   string    `gorm:"column:office_code;type:varchar(150);not null" json:"office_code"`
	ObjectCode   string    `gorm:"column:object_code;type:varchar(10);default:USER_OFFICE" json:"object_code"`
	TimezoneSet  string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (UserOffice) TableName() string {
	return "hrms_user_office"
}
