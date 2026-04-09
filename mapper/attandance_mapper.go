package mappers

import (
	"hrms_go/dto/attandance"
	"hrms_go/models"
	"time"
)

// adjust format based on your input (IMPORTANT)
const (
	dateFormat     = "2006-01-02"
	timeFormat     = "15:04:05"
	dateTimeFormat = "2006-01-02 15:04:05"
)

func ToAttendanceModel(dto attandance.PostAttandanceDto) (models.Attendance, error) {
	var result models.Attendance
	var err error

	// parse date
	var parsedDate time.Time
	if dto.Date != "" {
		parsedDate, err = time.Parse(dateFormat, dto.Date)
		if err != nil {
			return result, err
		}
	}

	// parse server timestamp
	var parsedServerTime time.Time
	if dto.ServerTimestamp != "" {
		parsedServerTime, err = time.Parse(dateTimeFormat, dto.ServerTimestamp)
		if err != nil {
			return result, err
		}
	}

	// parse created at
	var createdAt time.Time
	if dto.CreatedAt != "" {
		createdAt, err = time.Parse(dateTimeFormat, dto.CreatedAt)
		if err != nil {
			return result, err
		}
	}

	// parse updated at (nullable)
	var updatedAt *time.Time
	if dto.UpdatedAt != "" {
		t, err := time.Parse(dateTimeFormat, dto.UpdatedAt)
		if err != nil {
			return result, err
		}
		updatedAt = &t
	}

	result = models.Attendance{
		AttendanceId: dto.AttendanceId,
		UserId:       dto.UserId,
		DeviceId:     dto.DeviceId,

		CheckType:        dto.CheckType,
		CheckDescription: dto.CheckDescription,

		ShiftCode:          dto.ShiftCode,
		ShiftDurationHours: dto.ShiftDurationHours,

		Date:            parsedDate,
		Time:            dto.Time,
		ServerTimestamp: parsedServerTime,

		LocationCode: dto.LocationCode,
		LocationName: dto.LocationName,

		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,

		GpsAccuracy: dto.GpsAccuracy,

		IsMockLocation: dto.IsMockLocation,
		Activity:       dto.Activity,

		PhotoUrl: dto.PhotoUrl,
		Notes:    dto.Notes,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return result, nil
}