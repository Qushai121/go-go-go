package models

import (
	"time"

	"github.com/google/uuid"
)

type OfficePermit struct {
	OfficePermitId   uuid.UUID  `gorm:"column:office_permit_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"office_permit_id" form:"office_permit_id"`
	EmployeeNik      string     `gorm:"column:employee_nik;type:varchar(100);not null" json:"employee_nik" form:"employee_nik"`
	OfficePermitDate time.Time  `gorm:"column:office_permit_date;type:timestamp;not null" json:"office_permit_date" form:"office_permit_date"`
	Remarks          *string    `gorm:"column:remarks;type:text" json:"remarks" form:"remarks"`
	Status           string     `gorm:"column:status;type:varchar(20);default:P" json:"status" form:"status"`
	CurrentStep      int        `gorm:"column:current_step;type:int;default:1" json:"current_step" form:"current_step"`
	ApprovalHeaderId *uuid.UUID `gorm:"column:approvalheader_id;type:uuid" json:"approvalheader_id" form:"approvalheader_id"`
	ObjectCode       string     `gorm:"column:object_code;type:varchar(50);default:LEAVE_HISTORY" json:"object_code" form:"object_code"`
	CreatedBy        string     `gorm:"column:created_by;type:varchar(100);default:System" json:"created_by" form:"created_by"`
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedBy        *string    `gorm:"column:updated_by;type:varchar(100)" json:"updated_by" form:"updated_by"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at" form:"updated_at"`
}

func (OfficePermit) TableName() string {
	return "hrms_office_permit"
}
