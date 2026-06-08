package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Title struct {
	TitleId        uuid.UUID `gorm:"column:title_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"title_id"`
	CompanyCode    string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	DivisionCode   string    `gorm:"column:division_code;type:varchar(150);not null" json:"division_code"`
	DepartmentCode string    `gorm:"column:department_code;type:varchar(150);not null" json:"department_code"`
	TitleCode      string    `gorm:"column:title_code;type:varchar(150);not null" json:"title_code"`
	TitleName      string    `gorm:"column:title_name;type:varchar(150);not null" json:"title_name"`
	ObjectCode     string    `gorm:"column:object_code;type:varchar(10);default:TITLE" json:"object_code"`
	TimezoneSet    string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Title) TableName() string {
	return "hrms_title"
}
