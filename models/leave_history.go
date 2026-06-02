package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveHistory struct {
	LeaveHistoryId uuid.UUID `gorm:"column:leave_history_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"leave_history_id" form:"leave_history_id"`

	EmployeeNik string    `gorm:"column:employee_nik;type:varchar(100);not null" json:"employee_nik" form:"employee_nik"`
	LeaveType   string    `gorm:"column:leave_type;type:varchar(50);not null" json:"leave_type" form:"leave_type"`
	LeaveStart  time.Time `gorm:"column:leave_start;type:timestamp" json:"leave_start" form:"leave_start"`
	LeaveEnd    time.Time `gorm:"column:leave_end;type:timestamp" json:"leave_end" form:"leave_end"`

	TotalDays int     `gorm:"column:total_days;type:int" json:"total_days" form:"total_days"`
	Remarks   *string `gorm:"column:remarks;type:text" json:"remarks" form:"remarks"`
	Location  *string `gorm:"column:location;type:varchar(255)" json:"location" form:"location"`

	LeaveYear int    `gorm:"column:leave_year;type:int" json:"leave_year" form:"leave_year"`
	Status    string `gorm:"column:status;type:varchar(10)" json:"status" form:"status"`

	CurrentStep      int        `gorm:"column:current_step;type:int" json:"current_step" form:"current_step"`
	ApprovalHeaderId *uuid.UUID `gorm:"column:approvalheader_id;type:uuid" json:"approvalheader_id" form:"approvalheader_id"`

	ObjectCode string `gorm:"column:object_code;type:varchar(100)" json:"object_code" form:"object_code"`

	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp" json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at" form:"updated_at"`

	CreatedBy string  `gorm:"column:created_by;type:varchar(100)" json:"created_by" form:"created_by"`
	UpdatedBy *string `gorm:"column:updated_by;type:varchar(100)" json:"updated_by" form:"updated_by"`
}

func (LeaveHistory) TableName() string {
	return "hrms_leave_history"
}
