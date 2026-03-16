package models

import (
	"time"

	"github.com/google/uuid"
)

type Param struct {
	ParamId      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"param_id"`
	ParamGroupId uuid.UUID `gorm:"type:uuid;not null" json:"paramgroup_id"`

	ParamCode string `gorm:"type:varchar(150);not null" json:"param_code"`
	ParamName string `gorm:"type:varchar(150);not null" json:"param_name"`

	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	CreatedBy string  `gorm:"type:varchar(100);default:System" json:"created_by"`
	UpdatedBy *string `gorm:"type:varchar(100)" json:"updated_by"`

	ParamGroup ParamGroup `gorm:"foreignKey:ParamGroupId;references:ParamGroupId"`
}

func (Param) TableName() string {
	return "hrms_par"
}