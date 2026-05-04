package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type UserCompany struct {
	UserCompanyId uuid.UUID `gorm:"column:user_company_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"user_company_id"`
	EmployeeNIK   string    `gorm:"column:employee_nik;type:varchar(50);not null" json:"employee_nik"`
	CompanyCode   string    `gorm:"column:company_code;type:varchar(150);not null" json:"company_code"`
	ObjectCode    string    `gorm:"column:object_code;type:varchar(10);default:USER_COMPANY" json:"object_code"`
	TimezoneSet   string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (UserCompany) TableName() string {
	return "hrms_user_company"
}
