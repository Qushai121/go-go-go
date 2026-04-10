package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ApprovalDetail struct {
	ApprovalDetailId uuid.UUID `gorm:"column:approval_detail_id;type:uuid;default:uuid_generate_v4();primaryKey"`

	ApprovalHeaderId uuid.UUID `gorm:"column:approval_header_id;type:uuid"`

	ApprovalStatus string `gorm:"column:approval_status;type:varchar(50)"`

	ApproverBy uuid.UUID `gorm:"column:approver_by;type:uuid"`
	Approver       User      `gorm:"foreignKey:ApproverBy;references:UserId"`

	base.AuditFields
}

func (ApprovalDetail) TableName() string {
	return "hrms_approval_detail"
}