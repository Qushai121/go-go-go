package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ApprovalTemplateHeader struct {
	ApprovalTemplateHeaderId uuid.UUID `gorm:"column:approval_template_header_id;type:uuid;default:uuid_generate_v4();primaryKey"`

	TemplateType string `gorm:"column:template_type;type:varchar(100)"`

	base.AuditFields
}

func (ApprovalTemplateHeader) TableName() string {
	return "hrms_approval_template_header"
}