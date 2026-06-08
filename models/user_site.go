package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type UserSite struct {
	UserSiteId  uuid.UUID `gorm:"column:user_site_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"user_site_id"`
	CompanyCode string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	EmployeeNIK string    `gorm:"column:employee_nik;type:varchar(50);not null" json:"employee_nik"`
	SiteType    string    `gorm:"column:site_type;type:varchar(50);not null" json:"site_type"`
	SiteCode    string    `gorm:"column:site_code;type:varchar(255);not null" json:"site_code"`
	ObjectCode  string    `gorm:"column:object_code;type:varchar(10);default:USER_SITE" json:"object_code"`
	TimezoneSet string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	SiteId        *uuid.UUID `gorm:"->;column:site_id;-:migration" json:"site_id,omitempty"`
	SiteName      string     `gorm:"->;column:site_name;-:migration" json:"site_name,omitempty"`
	SitePhone     string     `gorm:"->;column:site_phone;-:migration" json:"site_phone,omitempty"`
	SiteAddress   string     `gorm:"->;column:site_address;-:migration" json:"site_address,omitempty"`
	SiteLatitude  string     `gorm:"->;column:site_latitude;-:migration" json:"site_latitude,omitempty"`
	SiteLongitude string     `gorm:"->;column:site_longitude;-:migration" json:"site_longitude,omitempty"`
	MaxRadius     int        `gorm:"->;column:max_radius;-:migration" json:"max_radius,omitempty"`

	base.AuditFields
}

func (UserSite) TableName() string {
	return "hrms_user_site"
}
