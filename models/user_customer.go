package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type UserCustomer struct {
	UserCustomerId uuid.UUID `gorm:"column:user_customer_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"user_customer_id"`
	EmployeeNIK    string    `gorm:"column:employee_nik;type:varchar(50);not null" json:"employee_nik"`
	CustomerCode   string    `gorm:"column:customer_code;type:varchar(150);not null" json:"customer_code"`
	ObjectCode     string    `gorm:"column:object_code;type:varchar(10);default:USER_CUSTOMER" json:"object_code"`
	TimezoneSet    string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (UserCustomer) TableName() string {
	return "hrms_user_customer"
}
