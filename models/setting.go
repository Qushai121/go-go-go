package models

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	SettingId    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"setting_id"`
	SettingName  string    `gorm:"type:varchar(150);not null;unique" json:"setting_name"`
	SettingValue string    `gorm:"type:varchar(150);not null" json:"setting_value"`

	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	CreatedBy string  `gorm:"type:varchar(100);default:System" json:"created_by"`
	UpdatedBy *string `gorm:"type:varchar(100)" json:"updated_by"`
}

func (Setting) TableName() string {
	return "hrms_setting"
}