package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Office struct {
	OfficeId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"office_id"`

	CompanyCode string `gorm:"type:varchar(50);not null" json:"company_code"`
	BranchCode  string `gorm:"type:varchar(50);not null" json:"branch_code"`
	OfficeCode  string `gorm:"type:varchar(255);not null" json:"office_code"`
	OfficeName  string `gorm:"type:varchar(100);not null" json:"office_name"`

	OfficePhone   string `gorm:"type:varchar(13)" json:"office_phone"`
	OfficeAddress string `gorm:"type:text" json:"office_address"`

	OfficeLatitude  string `gorm:"type:varchar(255)" json:"office_latitude"`
	OfficeLongitude string `gorm:"type:varchar(255)" json:"office_longitude"`

	MaxRadius int `gorm:"default:5" json:"max_radius"`

	ObjectCode string `gorm:"type:varchar(10);default:OFFICE" json:"object_code"`

	TimezoneSet string `gorm:"type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Office) TableName() string {
	return "hrms_office"
}
