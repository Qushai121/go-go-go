package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveHistory struct {
	LeaveHistoryId uuid.UUID `gorm:"column:leave_history_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"leave_history_id" form:"leave_history_id"`

	EmployeeNik string     `gorm:"column:employee_nik;type:varchar(100);not null" json:"employee_nik" form:"employee_nik"`
	LeaveTypeId *uuid.UUID `gorm:"column:leave_type_id;type:uuid" json:"leave_type_id" form:"leave_type_id"`
	LeaveType   string     `gorm:"column:leave_type;type:varchar(50);not null" json:"leave_type" form:"leave_type"`
	LeaveDate   time.Time  `gorm:"column:leave_date;type:date" json:"leave_date" form:"leave_date"`
	LeaveStart  time.Time  `gorm:"column:leave_start;type:timestamp" json:"leave_start" form:"leave_start"`
	LeaveEnd    time.Time  `gorm:"column:leave_end;type:timestamp" json:"leave_end" form:"leave_end"`

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

type LeaveHistoryFilter struct {
	LeaveHistoryId string `json:"leave_history_id" query:"leave_history_id" form:"leave_history_id"`

	EmployeeNik string `json:"employee_nik" query:"employee_nik" form:"employee_nik"`
	LeaveTypeId string `json:"leave_type_id" query:"leave_type_id" form:"leave_type_id"`
	LeaveType   string `json:"leave_type" query:"leave_type" form:"leave_type"`

	// Filter overlap cuti berdasarkan range tanggal
	StartDate string `json:"start_date" query:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" query:"end_date" form:"end_date"`

	// Filter tanggal exact dari kolom leave_date, leave_start, dan leave_end
	LeaveDate  string `json:"leave_date" query:"leave_date" form:"leave_date"`
	LeaveStart string `json:"leave_start" query:"leave_start" form:"leave_start"`
	LeaveEnd   string `json:"leave_end" query:"leave_end" form:"leave_end"`

	TotalDays string `json:"total_days" query:"total_days" form:"total_days"`
	Remarks   string `json:"remarks" query:"remarks" form:"remarks"`
	Location  string `json:"location" query:"location" form:"location"`

	LeaveYear string `json:"leave_year" query:"leave_year" form:"leave_year"`
	Status    string `json:"status" query:"status" form:"status"`

	CurrentStep      string `json:"current_step" query:"current_step" form:"current_step"`
	ApprovalHeaderId string `json:"approvalheader_id" query:"approvalheader_id" form:"approvalheader_id"`

	ObjectCode string `json:"object_code" query:"object_code" form:"object_code"`

	CreatedAt string `json:"created_at" query:"created_at" form:"created_at"`
	UpdatedAt string `json:"updated_at" query:"updated_at" form:"updated_at"`

	CreatedBy string `json:"created_by" query:"created_by" form:"created_by"`
	UpdatedBy string `json:"updated_by" query:"updated_by" form:"updated_by"`

	// Opsional range audit date
	CreatedAtStart string `json:"created_at_start" query:"created_at_start" form:"created_at_start"`
	CreatedAtEnd   string `json:"created_at_end" query:"created_at_end" form:"created_at_end"`
	UpdatedAtStart string `json:"updated_at_start" query:"updated_at_start" form:"updated_at_start"`
	UpdatedAtEnd   string `json:"updated_at_end" query:"updated_at_end" form:"updated_at_end"`
}

func (LeaveHistory) TableName() string {
	return "hrms_leave_history"
}
