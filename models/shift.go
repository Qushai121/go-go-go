package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Shift struct {
	ShiftId      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"shift_id"`
	ShiftCode     string    `json:"shift_code"`
	ShiftName     string    `json:"shift_name"`
	ShiftDuration int    `json:"shift_duration"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	GracePeriod    int       `json:"grace_period"`
	base.AuditFields
}


func (s Shift) TableName() string {
	return "hrms_shift"
}
