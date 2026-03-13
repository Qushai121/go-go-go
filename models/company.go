package models

import "hrms_go/models/base"

type Company struct {
	CompanyId   string `gorm:"primaryKey" json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`

	base.AuditFields
}

func (c Company) TableName() string {
	return "hrms_company"
}