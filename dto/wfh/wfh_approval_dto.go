package wfh

import (
	"time"

	"github.com/google/uuid"
)

type WFHApproval struct {
	WFHId     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"wfh_id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" form:"user_id"`
	Remarks   string    `json:"remarks"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	
	ApprovalHeaderId *uuid.UUID `json:"approval_header_id"`
	FinalStatus  string `json:"final_status"`
}