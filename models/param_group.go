package models

import (
	"time"

	"github.com/google/uuid"
)

type ParamGroup struct {
	ParamGroupId   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"paramgroup_id"`
	ParamGroupCode string    `gorm:"type:varchar(50);unique;not null" json:"paramgroup_code"`
	ParamGroupName string    `gorm:"type:varchar(150);not null" json:"paramgroup_name"`

	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedBy string     `gorm:"type:varchar(100);default:System" json:"created_by"`
	UpdatedBy *string    `gorm:"type:varchar(100)" json:"updated_by"`
}

func (ParamGroup) TableName() string {
	return "hrms_parg"
}
