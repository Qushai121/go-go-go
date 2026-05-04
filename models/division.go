package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Division struct {
	DivisionId   uuid.UUID `gorm:"column:division_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"division_id"`
	CompanyCode  string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	BranchCode   string    `gorm:"column:branch_code;type:varchar(150);not null" json:"branch_code"`
	OfficeCode   string    `gorm:"column:office_code;type:varchar(150);not null" json:"office_code"`
	DivisionCode string    `gorm:"column:division_code;type:varchar(150);not null" json:"division_code"`
	DivisionName string    `gorm:"column:division_name;type:varchar(150);not null" json:"division_name"`
	ObjectCode   string    `gorm:"column:object_code;type:varchar(10);default:DIVISION" json:"object_code"`
	TimezoneSet  string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Division) TableName() string {
	return "hrms_division"
}
