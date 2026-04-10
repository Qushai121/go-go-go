package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ApprovalTemplateDetail struct {
	ApprovalTemplateDetailId uuid.UUID `gorm:"column:approval_template_detail_id;type:uuid;default:uuid_generate_v4();primaryKey"`

	ApprovalTemplateHeaderId uuid.UUID               `gorm:"column:approval_template_header_id;type:uuid"`

	ApproverBy uuid.UUID `gorm:"column:approver_by;type:uuid"`
	Approver   User      `gorm:"foreignKey:ApproverBy;references:UserId"`

	SequenceNumber int `gorm:"column:sequence_number"`

	base.AuditFields
}

func (ApprovalTemplateDetail) TableName() string {
	return "hrms_approval_template_detail"
}