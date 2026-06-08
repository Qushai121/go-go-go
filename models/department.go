package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Department struct {
	DepartmentId   uuid.UUID `gorm:"column:department_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"department_id"`
	CompanyCode    string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	DivisionCode   string    `gorm:"column:division_code;type:varchar(150);not null" json:"division_code"`
	DepartmentCode string    `gorm:"column:department_code;type:varchar(150);not null" json:"department_code"`
	DepartmentName string    `gorm:"column:department_name;type:varchar(150);not null" json:"department_name"`
	ObjectCode     string    `gorm:"column:object_code;type:varchar(10);default:DEPARTMENT" json:"object_code"`
	TimezoneSet    string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Department) TableName() string {
	return "hrms_department"
}
