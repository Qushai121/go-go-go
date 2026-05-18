package models

import (
	"hrms_go/models/base"
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	SubmissionId     uuid.UUID  `gorm:"column:submission_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"submission_id"`
	ReceiptId        uuid.UUID  `gorm:"column:receipt_id;type:uuid;not null" json:"receipt_id"`
	SubmissionNumber string     `gorm:"column:submission_number;type:varchar(100);not null" json:"submission_number"`
	SubmissionDate   time.Time  `gorm:"column:submission_date;not null" json:"submission_date"`
	Status           string     `gorm:"column:status;type:varchar(20);default:P" json:"status"`
	CurrentStep      int        `gorm:"column:current_step;default:1" json:"current_step"`
	ApprovalHeaderId *uuid.UUID `gorm:"column:approvalheader_id;type:uuid" json:"approvalheader_id"`
	ObjectCode       string     `gorm:"column:object_code;type:varchar(50);default:SUBMISSION" json:"object_code"`

	base.AuditFields
}

func (Submission) TableName() string {
	return "hrms_submission"
}
