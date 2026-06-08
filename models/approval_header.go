package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ApprovalHeader struct {
	ApprovalHeaderId uuid.UUID `gorm:"column:approval_header_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"approval_header_id"`

	ApprovalTemplateHeaderId uuid.UUID              `gorm:"column:approval_template_header_id;type:uuid"`
	ApprovalTemplateHeader   ApprovalTemplateHeader `gorm:"foreignKey:ApprovalTemplateHeaderId;references:ApprovalTemplateHeaderId"`

	ApprovalDocId uuid.UUID `gorm:"column:approval_doc_id;type:uuid"`

	RequesterBy uuid.UUID `gorm:"column:requester_by;type:uuid"`
	Requester   User      `gorm:"foreignKey:RequesterBy;references:UserId"`

	ApprovalDetails []ApprovalDetail `gorm:"foreignKey:ApprovalHeaderId;references:ApprovalHeaderId"`

	base.AuditFields
}

func (ApprovalHeader) TableName() string {
	return "hrms_approval_header"
}
