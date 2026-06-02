package models

import "time"

type UserShiftMapping struct {
	WeekStartDate time.Time `gorm:"column:week_start_date" json:"week_start_date"`
	WeekEndDate   time.Time `gorm:"column:week_end_date" json:"week_end_date"`

	ShiftDate time.Time `gorm:"column:shift_date" json:"shift_date"`

	EmployeeNIK string `gorm:"column:employee_nik" json:"employee_nik"`

	ShiftName        string `gorm:"column:shift_name" json:"shift_name"`
	IsActive         string `gorm:"column:is_active" json:"is_active"`
	WorkingHoursType string `gorm:"column:working_hours_type" json:"working_hours_type"`

	WeekdayID   int    `gorm:"column:weekday_id" json:"weekday_id"`
	WeekdayName string `gorm:"column:weekday_name" json:"weekday_name"`

	EventCode string `gorm:"column:event_code" json:"event_code"`
	EventName string `gorm:"column:event_name" json:"event_name"`

	StartTime string `gorm:"column:start_time" json:"start_time"`
	EndTime   string `gorm:"column:end_time" json:"end_time"`

	ToleranceBefore int `gorm:"column:tolerance_before" json:"tolerance_before"`
	TolaranceAfter  int `gorm:"column:tolarance_after" json:"tolarance_after"`
}
