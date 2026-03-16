package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Companies struct {
	CompaniesId   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"companies_id"`
	CompaniesCode string `json:"companies_code"`
	CompaniesName string `json:"companies_name"`

	base.AuditFields
}

func (c Companies) TableName() string {
	return "hrms_companies"
}