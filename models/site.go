package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Site struct {
	SiteId        uuid.UUID `gorm:"column:site_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"site_id"`
	CompanyCode   string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	SiteType      string    `gorm:"column:site_type;type:varchar(50);not null" json:"site_type"`
	SiteCode      string    `gorm:"column:site_code;type:varchar(255);not null" json:"site_code"`
	SiteName      string    `gorm:"column:site_name;type:varchar(100);not null" json:"site_name"`
	SitePhone     string    `gorm:"column:site_phone;type:varchar(13)" json:"site_phone"`
	SiteAddress   string    `gorm:"column:site_address;type:text" json:"site_address"`
	SiteLatitude  string    `gorm:"column:site_latitude;type:varchar(255)" json:"site_latitude"`
	SiteLongitude string    `gorm:"column:site_longitude;type:varchar(255)" json:"site_longitude"`
	MaxRadius     int       `gorm:"column:max_radius;default:5" json:"max_radius"`
	ObjectCode    string    `gorm:"column:object_code;type:varchar(10);default:SITE" json:"object_code"`
	TimezoneSet   string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Site) TableName() string {
	return "hrms_site"
}
