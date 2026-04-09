package attandance

import (
	"github.com/google/uuid"
)

type PostAttandanceDto struct {
	AttendanceId uuid.UUID `json:"attendance_id" form:"attendance_id"`

	UserId   uuid.UUID `json:"user_id" form:"user_id"`
	DeviceId string    `json:"device_id" form:"device_id"`

	CheckType        string `json:"check_type" form:"check_type"`
	CheckDescription string `json:"check_description" form:"check_description"`

	ShiftCode          string `json:"shift_code" form:"shift_code"`
	ShiftDurationHours int    `json:"shift_duration_hours" form:"shift_duration_hours"`

	Date            string `json:"date" form:"date"`
	Time            string    `json:"time" form:"time"`
	ServerTimestamp string `json:"server_timestamp" form:"server_timestamp"`

	LocationCode string `json:"location_code" form:"location_code"`
	LocationName string `json:"location_name" form:"location_name"`

	Latitude  float64 `json:"latitude" form:"latitude"`
	Longitude float64 `json:"longitude" form:"longitude"`

	GpsAccuracy float64 `json:"gps_accuracy" form:"gps_accuracy"`

	IsMockLocation bool   `json:"is_mock_location" form:"is_mock_location"`
	Activity       string `json:"activity" form:"activity"`

	PhotoUrl string `json:"photo_url" form:"photo_url"`
	Notes    string `json:"notes" form:"notes"`

	CreatedAt string  `json:"created_at" form:"created_at"`
	UpdatedAt string `json:"updated_at" form:"updated_at"`
}