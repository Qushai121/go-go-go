package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	AttendanceId uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"attendance_id"`

	UserId   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	DeviceId string    `gorm:"type:varchar(150)" json:"device_id"`

	CheckType        string `gorm:"type:varchar(50);not null" json:"check_type"`
	CheckDescription string `gorm:"type:varchar(150)" json:"check_description"`

	ShiftCode          string `gorm:"type:varchar(50)" json:"shift_code"`
	ShiftDurationHours int    `json:"shift_duration_hours"`

	Date            time.Time `gorm:"type:date;not null" json:"date"`
	Time            string    `gorm:"type:time;not null" json:"time"`
	ServerTimestamp time.Time `json:"server_timestamp"`

	LocationCode string `gorm:"type:varchar(100)" json:"location_code"`
	LocationName string `gorm:"type:varchar(150)" json:"location_name"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	GpsAccuracy float64 `json:"gps_accuracy"`

	IsMockLocation bool   `json:"is_mock_location"`
	Activity       string `gorm:"type:varchar(100)" json:"activity"`

	PhotoUrl string `gorm:"type:text" json:"photo_url"`
	Notes    string `gorm:"type:text" json:"notes"`

	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserId;references:UserId"`
}

func (Attendance) TableName() string {
	return "hrms_attendance"
}