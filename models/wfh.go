package models

import (
	"hrms_go/models/base"
	"time"

	"github.com/google/uuid"
)

type WFH struct {
	WFHId 		uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"wfh_id"`
	UserId   	uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" form:"user_id"`
	Remarks       string    `json:"remarks"`
	StartTime     time.Time    `json:"start_time"`
	EndTime       time.Time    `json:"end_time"`
	// User User `gorm:"foreignKey:UserId;references:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	base.AuditFields
}

func (WFH) TableName() string {
	return "hrms_wfh"
}
