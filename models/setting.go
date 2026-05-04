package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Setting struct {
	SettingId    uuid.UUID `gorm:"column:setting_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"setting_id"`
	CompanyCode  string    `gorm:"column:company_code;type:varchar(100);not null" json:"company_code"`
	SettingCode  string    `gorm:"column:setting_code;type:varchar(100);not null" json:"setting_code"`
	SettingName  string    `gorm:"column:setting_name;type:varchar(150);not null" json:"setting_name"`
	SettingValue string    `gorm:"column:setting_value;type:varchar(150);not null" json:"setting_value"`

	base.AuditFields
}

func (Setting) TableName() string {
	return "hrms_setting"
}
