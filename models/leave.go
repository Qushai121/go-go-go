package models

import (
	"hrms_go/models/base"
	"time"

	"github.com/google/uuid"
)

type Leave struct {
	LeaveID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"leave_id"`
    RequestNumber     string    `json:"request_number"`
    EmployeeName      string    `json:"employee_name"`
    LeaveType         string    `json:"leave_type"`
    StartDate         time.Time `json:"start_date"`
    EndDate           time.Time `json:"end_date"`
    ReqDate           time.Time `json:"req_date"`
    Status            string    `json:"status"`
    Message           string    `json:"message"`
    CancellationReason string   `json:"cancellation_reason"`
    LeaveBalance      float64   `json:"leave_balance"`

	base.AuditFields
}

func (Leave) TableName() string {
	return "hrms_leave"
}