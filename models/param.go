package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Param struct {
	ParamId      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"param_id"`
	ParamGroupId uuid.UUID `gorm:"column:paramgroup_id;type:uuid;not null" json:"paramgroup_id"`
	CompanyCode  string    `gorm:"column:company_code;type:varchar(100);not null" json:"company_code"`

	ParamCode string `gorm:"type:varchar(150);not null" json:"param_code"`
	ParamName string `gorm:"type:varchar(150);not null" json:"param_name"`

	ParamGroup ParamGroup `gorm:"foreignKey:ParamGroupId;references:ParamGroupId"`

	base.AuditFields
}

func (Param) TableName() string {
	return "hrms_par"
}
