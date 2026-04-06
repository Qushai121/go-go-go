package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Customer struct {
	CustomerId      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"customer_id"`
	LocationCode    string    `json:"location_code"`
	LocationName    string    `json:"location_name"`
	Address         string    `json:"address"`
	TargetLatitude  string    `json:"target_latitude"`
	TargetLongitude string    `json:"target_longitude"`
	RadiusMeter     int       `json:"radius_meter"`

	base.AuditFields
}

func (s Customer) TableName() string {
	return "hrms_customer"
}
