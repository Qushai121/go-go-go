package models

import (
	"hrms_go/models/base"
	"time"
)

type Shift struct {
	ShiftId       string    `gorm:"primaryKey" json:"shift_id"`
	ShiftCode     string    `json:"shift_code"`
	ShiftName     string    `json:"shift_name"`
	ShiftDuration string    `json:"shift_duration"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	GracePeriod    int       `json:"grace_period"`
	
	base.AuditFields
}


func (s Shift) TableName() string {
	return "hrms_shift"
}
