package models

import "hrms_go/models/base"

type Customer struct {
	CustomerId      string `gorm:"primaryKey" json:"customer_id"`
	LocationCode    string `json:"location_code"`
	LocationName    string `json:"location_name"`
	Address         string `json:"address"`
	TargetLatitude  string `json:"target_latitude"`
	TargetLongitude string `json:"target_longitude"`
	RadiusMeter     int    `json:"radius_meter"`

	base.AuditFields
}

func (s Customer) TableName() string {
	return "hrms_customer"
}
