package models

import (
	"hrms_go/constant"
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type ApprovalDetail struct {
	ApprovalDetailId uuid.UUID `gorm:"column:approval_detail_id;type:uuid;default:uuid_generate_v4();primaryKey"`

	ApprovalHeaderId uuid.UUID `gorm:"column:approval_header_id;type:uuid"`

	ApprovalStatus constant.ApprovalStatus `gorm:"column:approval_status;type:varchar(50)"`

	ApproverBy uuid.UUID `gorm:"column:approver_by;type:uuid"`
	Approver   User      `gorm:"foreignKey:ApproverBy;references:UserId"`

	Remark string `gorm:"column:remark;type:varchar(255)"`

	base.AuditFields
}

func (ApprovalDetail) TableName() string {
	return "hrms_approval_detail"
}
