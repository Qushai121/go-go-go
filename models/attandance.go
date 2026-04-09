package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	AttendanceId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"attendance_id" form:"attendance_id"`

	UserId   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" form:"user_id"`
	DeviceId string    `gorm:"type:varchar(150);not null" json:"device_id" form:"device_id"`

	CheckType        string `gorm:"type:varchar(50);not null;index" json:"check_type" form:"check_type"`
	CheckDescription string `gorm:"type:varchar(150)" json:"check_description" form:"check_description"`

	ShiftCode          string `gorm:"type:varchar(50);index" json:"shift_code" form:"shift_code"`
	ShiftDurationHours int    `gorm:"type:int;default:0" json:"shift_duration_hours" form:"shift_duration_hours"`

	Date            time.Time `gorm:"type:date;not null;index" json:"date" form:"date"`
	Time            string    `gorm:"type:time;not null" json:"time" form:"time"`
	ServerTimestamp time.Time `gorm:"autoCreateTime" json:"server_timestamp" form:"server_timestamp"`

	LocationCode string `gorm:"type:varchar(100);index" json:"location_code" form:"location_code"`
	LocationName string `gorm:"type:varchar(150)" json:"location_name" form:"location_name"`

	Latitude  float64 `gorm:"type:decimal(10,7)" json:"latitude" form:"latitude"`
	Longitude float64 `gorm:"type:decimal(10,7)" json:"longitude" form:"longitude"`

	GpsAccuracy float64 `gorm:"type:decimal(10,2)" json:"gps_accuracy" form:"gps_accuracy"`

	IsMockLocation bool   `gorm:"default:false" json:"is_mock_location" form:"is_mock_location"`
	Activity       string `gorm:"type:varchar(100)" json:"activity" form:"activity"`

	PhotoUrl string `gorm:"type:text" json:"photo_url" form:"photo_url"`
	Notes    string `gorm:"type:text" json:"notes" form:"notes"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at" form:"updated_at"`

	User User `gorm:"foreignKey:UserId;references:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Attendance) TableName() string {
	return "hrms_attendance"
}