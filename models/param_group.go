package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ParamGroup struct {
	ParamGroupId   uuid.UUID `gorm:"column:paramgroup_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"paramgroup_id"`
	CompanyCode    string    `gorm:"column:company_code;type:varchar(100);not null" json:"company_code"`
	ParamGroupCode string    `gorm:"column:paramgroup_code;type:varchar(50);not null" json:"paramgroup_code"`
	ParamGroupName string    `gorm:"column:paramgroup_name;type:varchar(150);not null" json:"paramgroup_name"`

	base.AuditFields
}

func (ParamGroup) TableName() string {
	return "hrms_parg"
}
