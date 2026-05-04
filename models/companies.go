package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Companies struct {
	CompanyId   uuid.UUID `gorm:"column:company_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"company_id"`
	CompanyCode string    `gorm:"column:company_code;type:varchar;not null;unique" json:"company_code"`
	CompanyName string    `gorm:"column:company_name;type:varchar;not null" json:"company_name"`
	ObjectCode  string    `gorm:"column:object_code;type:varchar(10);default:COMPANY" json:"object_code"`

	base.AuditFields
}

func (c Companies) TableName() string {
	return "hrms_company"
}
