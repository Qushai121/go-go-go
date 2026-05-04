package mappers

import (
	"fmt"
	"hrms_go/dto/attandance"
	"hrms_go/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	dateTimeLayout       = "2006-01-02 15:04:05"
	dateTimeLayoutISO    = time.RFC3339
	dateTimeLayoutNoZone = "2006-01-02T15:04:05"
)

func ToAttendanceModel(dto attandance.PostAttandanceDto) (models.Attendance, error) {
	var result models.Attendance

	userId, err := parseUUIDField(dto.UserId, "user_id", true)
	if err != nil {
		return result, err
	}

	attendanceId, err := parseUUIDField(dto.AttendanceId, "attendance_id", false)
	if err != nil {
		return result, err
	}
	if attendanceId == uuid.Nil {
		attendanceId = uuid.New()
	}

	logTime, err := parseTimeField(dto.LogTime, "logtime", true)
	if err != nil {
		return result, err
	}

	maxRadius, err := parseOptionalIntField(dto.MaxRadius, "max_radius")
	if err != nil {
		return result, err
	}

	expandRadius, err := parseOptionalIntField(dto.ExpandRadius, "expand_radius")
	if err != nil {
		return result, err
	}

	result = models.Attendance{
		AttendanceId: attendanceId,
		UserId:       userId,
		CompanyCode:  strings.TrimSpace(dto.CompanyCode),
		OfficeCode:   strings.TrimSpace(dto.OfficeCode),
		LogTime:      logTime,
		FunctionNo:   dto.FunctionNo,
		ActivityType: toOptionalString(dto.ActivityType),
		Latitude:     toOptionalString(dto.Latitude),
		Longitude:    toOptionalString(dto.Longitude),
		PresentaseKemiripan: toOptionalString(dto.PresentaseKemiripan),
		ImagePath:    strings.TrimSpace(dto.ImagePath),
		IsOffline:    toOptionalString(dto.IsOffline),
		Distance:     toOptionalString(dto.Distance),
		Platforms:    toOptionalString(dto.Platforms),
		MaxRadius:    maxRadius,
		ExpandRadius: expandRadius,
		ObjectCode:   strings.TrimSpace(dto.ObjectCode),
		CreatedBy:    strings.TrimSpace(dto.CreatedBy),
		UpdatedBy:    toOptionalStringPtr(dto.UpdatedBy),
	}

	if result.ObjectCode == "" {
		result.ObjectCode = "ATTENDANCE"
	}
	if result.CreatedBy == "" {
		result.CreatedBy = "system"
	}

	createdAt, err := parseTimeField(dto.CreatedAt, "created_at", false)
	if err != nil {
		return result, err
	}
	if !createdAt.IsZero() {
		result.CreatedAt = createdAt
	}

	updatedAt, err := parseTimeField(dto.UpdatedAt, "updated_at", false)
	if err != nil {
		return result, err
	}
	if !updatedAt.IsZero() {
		result.UpdatedAt = &updatedAt
	}

	return result, nil
}

func parseUUIDField(value string, fieldName string, required bool) (uuid.UUID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		if required {
			return uuid.Nil, fmt.Errorf("%s is required", fieldName)
		}
		return uuid.Nil, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", fieldName)
	}

	return parsed, nil
}

func parseTimeField(value string, fieldName string, required bool) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		if required {
			return time.Time{}, fmt.Errorf("%s is required", fieldName)
		}
		return time.Time{}, nil
	}

	layouts := []string{
		dateTimeLayout,
		dateTimeLayoutISO,
		dateTimeLayoutNoZone,
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid %s", fieldName)
}

func parseOptionalIntField(value string, fieldName string) (*int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s", fieldName)
	}

	return &parsed, nil
}

func toOptionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return &value
}

func toOptionalStringPtr(value string) *string {
	return toOptionalString(value)
}
