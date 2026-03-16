package models

import (
	"time"

	"github.com/google/uuid"
)

type Office struct {
	OfficeId uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"office_id"`

	CompanyCode string `gorm:"type:varchar(50);not null" json:"company_code"`
	OfficeCode  string `gorm:"type:varchar(255);not null" json:"office_code"`
	OfficeName  string `gorm:"type:varchar(100);not null" json:"office_name"`

	OfficePhone   string `gorm:"type:varchar(13)" json:"office_phone"`
	OfficeAddress string `gorm:"type:text" json:"office_address"`

	OfficeProvince    string `gorm:"type:varchar(255)" json:"office_province"`
	OfficeCity        string `gorm:"type:varchar(255)" json:"office_city"`
	OfficeSubdistrict string `gorm:"type:varchar(255)" json:"office_subdistrict"`
	OfficeWard        string `gorm:"type:varchar(255)" json:"office_ward"`

	OfficeLatitude  string `gorm:"type:varchar(255)" json:"office_latitude"`
	OfficeLongitude string `gorm:"type:varchar(255)" json:"office_longitude"`

	MaxRadius int `gorm:"default:5" json:"max_radius"`

	ObjectCode string `gorm:"type:varchar(10);default:OFFICE" json:"object_code"`

	TimezoneSet string `gorm:"type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	CreatedUser string    `gorm:"type:varchar(50);default:system" json:"created_user"`
	CreatedDate time.Time `gorm:"default:now()" json:"created_date"`

	UpdatedUser *string    `gorm:"type:varchar(50)" json:"updated_user"`
	UpdatedDate *time.Time `json:"updated_date"`

	CurrentUtcOffset string `gorm:"type:varchar(10);default:+07:00" json:"current_utc_offset"`

	OfficeCodeSunfish string `gorm:"type:varchar(50)" json:"office_code_sunfish"`
	OfficeCodeHa      string `gorm:"type:varchar(50)" json:"office_code_ha"`

	OfficePic string `gorm:"type:varchar(100)" json:"office_pic"`
}

func (Office) TableName() string {
	return "hrms_office"
}