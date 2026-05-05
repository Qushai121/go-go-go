package utils

import (
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func normalizeResponseData(data interface{}) interface{} {
	if data == nil {
		return fiber.Map{}
	}

	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return fiber.Map{}
		}
		return normalizeResponseData(value.Elem().Interface())
	case reflect.Slice:
		if value.IsNil() {
			return []any{}
		}
		return data
	case reflect.Array:
		return data
	case reflect.Map:
		if value.IsNil() {
			return fiber.Map{}
		}
		return data
	default:
		return data
	}
}

func extractResponseMessage(data interface{}) string {
	if data == nil {
		return "Success"
	}

	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "Success"
		}
		value = value.Elem()
	}

	if value.Kind() == reflect.String {
		return value.String()
	}

	return "Success"
}

func normalizeSuccessData(data interface{}) interface{} {
	if data == nil {
		return fiber.Map{}
	}

	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return fiber.Map{}
		}
		value = value.Elem()
	}

	if value.Kind() == reflect.String {
		return fiber.Map{}
	}

	return normalizeResponseData(data)
}

func isPaginatedResponse(data interface{}) bool {
	if data == nil {
		return false
	}

	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return false
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return false
	}

	_, hasData := value.Type().FieldByName("Data")
	_, hasTotalRecord := value.Type().FieldByName("TotalRecord")
	_, hasTotalPage := value.Type().FieldByName("TotalPage")

	return hasData && hasTotalRecord && hasTotalPage
}

func getPaginatedFieldValue(v reflect.Value, fieldName string) interface{} {
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil
		}
		field = field.Elem()
	}

	if !field.IsValid() || !field.CanInterface() {
		return nil
	}

	return field.Interface()
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}

func Success(ctx fiber.Ctx, data interface{}) error {
	if isPaginatedResponse(data) {
		value := reflect.ValueOf(data)
		for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		pageNo := parsePositiveInt(ctx.Get("pageNo"), 1)
		if pageNo == 1 {
			pageNo = parsePositiveInt(ctx.Query("pageNo"), 1)
		}
		if pageNo == 1 {
			pageNo = parsePositiveInt(ctx.Query("page"), 1)
		}

		pageSize := parsePositiveInt(ctx.Get("pageSize"), 10)
		if pageSize == 10 {
			pageSize = parsePositiveInt(ctx.Query("pageSize"), 10)
		}
		if pageSize == 10 {
			pageSize = parsePositiveInt(ctx.Query("per_page"), 10)
		}

		return ctx.Status(200).JSON(fiber.Map{
			"status":      200,
			"error":       nil,
			"message":     "Success",
			"data":        normalizeResponseData(getPaginatedFieldValue(value, "Data")),
			"pageNo":      pageNo,
			"pageSize":    pageSize,
			"pageTotal":   getPaginatedFieldValue(value, "TotalPage"),
			"totalRecord": getPaginatedFieldValue(value, "TotalRecord"),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": extractResponseMessage(data),
		"error":   nil,
		"data":    normalizeSuccessData(data),
	})
}

func Error(ctx fiber.Ctx, status int, err string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"status": status,
		"error":  err,
		"data":   fiber.Map{},
	})
}
