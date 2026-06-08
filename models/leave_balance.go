package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveBalance struct {
	LeaveBalanceId   uuid.UUID  `gorm:"column:leave_balance_id;type:uuid;primaryKey" json:"leave_balance_id"`
	EmployeeNik      string     `gorm:"column:employee_nik" json:"employee_nik"`
	LeaveTypeId      *uuid.UUID `gorm:"column:leave_type_id;type:uuid" json:"leave_type_id"`
	LeaveTypeCode    string     `gorm:"column:leave_type_code" json:"leave_type_code,omitempty"`
	LeaveTypeName    string     `gorm:"column:leave_type_name" json:"leave_type_name,omitempty"`
	LeavePeriodStart time.Time  `gorm:"column:leave_period_start" json:"leave_period_start"`
	LeavePeriodEnd   time.Time  `gorm:"column:leave_period_end" json:"leave_period_end"`
	TotalLeave       int        `gorm:"column:total_leave" json:"total_leave"`
	LeaveUsed        int        `gorm:"column:leave_used" json:"leave_used"`
	CarryForward     int        `gorm:"column:carry_forward" json:"carry_forward"`
	LeaveRemaining   int        `gorm:"column:leave_remaining" json:"leave_remaining"`
	ObjectCode       string     `gorm:"column:object_code" json:"object_code"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy        string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy        *string    `gorm:"column:updated_by" json:"updated_by"`
}

func (LeaveBalance) TableName() string {
	return "hrms_leave_balance"
}
