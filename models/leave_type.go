package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveType struct {
	LeaveTypeId          uuid.UUID  `gorm:"column:leave_type_id;type:uuid;primaryKey" json:"leave_type_id" form:"leave_type_id"`
	CompanyCode          string     `gorm:"column:company_code;type:varchar(50);not null" json:"company_code" form:"company_code"`
	LeaveTypeCode        string     `gorm:"column:leave_type_code;type:varchar(50);not null" json:"leave_type_code" form:"leave_type_code"`
	LeaveTypeName        string     `gorm:"column:leave_type_name;type:varchar(150);not null" json:"leave_type_name" form:"leave_type_name"`
	LeaveTypeDescription *string    `gorm:"column:leave_type_description;type:text" json:"leave_type_description" form:"leave_type_description"`
	IsActive             bool       `gorm:"column:is_active" json:"is_active" form:"is_active"`
	DeductLeaveBalance   bool       `gorm:"column:deduct_leave_balance" json:"deduct_leave_balance" form:"deduct_leave_balance"`
	DeductSalary         bool       `gorm:"column:deduct_salary" json:"deduct_salary" form:"deduct_salary"`
	IsPaidLeave          bool       `gorm:"column:is_paid_leave" json:"is_paid_leave" form:"is_paid_leave"`
	RequireAttachment    bool       `gorm:"column:require_attachment" json:"require_attachment" form:"require_attachment"`
	UseWorkingDay        bool       `gorm:"column:use_working_day" json:"use_working_day" form:"use_working_day"`
	AllowHalfDay         bool       `gorm:"column:allow_half_day" json:"allow_half_day" form:"allow_half_day"`
	AnnualQuota          int        `gorm:"column:annual_quota" json:"annual_quota" form:"annual_quota"`
	MaxDaysPerRequest    *int       `gorm:"column:max_days_per_request" json:"max_days_per_request" form:"max_days_per_request"`
	MinNoticeDays        int        `gorm:"column:min_notice_days" json:"min_notice_days" form:"min_notice_days"`
	SortOrder            int        `gorm:"column:sort_order" json:"sort_order" form:"sort_order"`
	ObjectCode           string     `gorm:"column:object_code;type:varchar(100)" json:"object_code" form:"object_code"`
	TimezoneSet          string     `gorm:"column:timezone_set;type:varchar(100)" json:"timezone_set" form:"timezone_set"`
	CreatedAt            time.Time  `gorm:"column:created_at;type:timestamp" json:"created_at" form:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp" json:"updated_at" form:"updated_at"`
	CreatedBy            string     `gorm:"column:created_by;type:varchar(100)" json:"created_by" form:"created_by"`
	UpdatedBy            *string    `gorm:"column:updated_by;type:varchar(100)" json:"updated_by" form:"updated_by"`
}

func (LeaveType) TableName() string {
	return "hrms_leave_type"
}

type LeaveTypeFilter struct {
	CompanyCode        string `json:"company_code" query:"company_code" form:"company_code"`
	Search             string `json:"search" query:"search" form:"search"`
	DeductLeaveBalance string `json:"deduct_leave_balance" query:"deduct_leave_balance" form:"deduct_leave_balance"`
	IsActive           string `json:"is_active" query:"is_active" form:"is_active"`
}
