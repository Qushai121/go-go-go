package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Customer struct {
	CustomerId        uuid.UUID `gorm:"column:customer_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"customer_id"`
	CustomerCode      string    `gorm:"column:customer_code;type:varchar(50);not null;unique" json:"customer_code"`
	CustomerName      string    `gorm:"column:customer_name;type:varchar(150);not null" json:"customer_name"`
	CustomerAddress   string    `gorm:"column:customer_address;type:text" json:"customer_address"`
	CustomerLatitude  string    `gorm:"column:customer_latitude;type:varchar(255);not null" json:"customer_latitude"`
	CustomerLongitude string    `gorm:"column:customer_longitude;type:varchar(255);not null" json:"customer_longitude"`
	MaxRadius         int       `gorm:"column:max_radius;default:5;not null" json:"max_radius"`
	ObjectCode        string    `gorm:"column:object_code;type:varchar(10);default:CUSTOMER" json:"object_code"`
	TimezoneSet       string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (s Customer) TableName() string {
	return "hrms_customer"
}
