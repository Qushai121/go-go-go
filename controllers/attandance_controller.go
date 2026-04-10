package controllers

import (
	"hrms_go/dto"
	"hrms_go/dto/attandance"
	mappers "hrms_go/mapper"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type AttendanceController struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceController(repo repositories.AttendanceRepository) *AttendanceController {
	return &AttendanceController{repo}
}

// Create Attendance godoc
// @Summary Create attendance
// @Description Create attendance with photo upload
// @Tags Attendance
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string false "User ID (UUID)"
// @Param device_id formData string false "Device ID"
// @Param activity formData string true "Activity" Enums(WFH, VISIT, OFFICE)
// @Param check_type formData string true "Check type" Enums(1,2)
// @Param check_description formData string false "Description"
// @Param shift_code formData string false "Shift code"
// @Param shift_duration_hours formData int false "Shift duration"
// @Param date formData string true "Date (YYYY-MM-DD)"
// @Param time formData string true "Time (HH:mm:ss)"
// @Param location_code formData string false "Location code"
// @Param location_name formData string false "Location name"
// @Param latitude formData number false "Latitude"
// @Param longitude formData number false "Longitude"
// @Param gps_accuracy formData number false "GPS accuracy"
// @Param is_mock_location formData boolean false "Mock location"
// @Param notes formData string false "Notes"
// @Param photo_url formData file true "Photo"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/attendance [post]
func (c *AttendanceController) Create(ctx fiber.Ctx) error {
	var request attandance.PostAttandanceDto

	if err := ctx.Bind().Body(&request); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := mappers.ToAttendanceModel(request); 
	if err != nil{
		return utils.Error(ctx, 400, err.Error())
	}

	file, err := ctx.FormFile("photo_url")
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	fileUrl, err := utils.SaveFileToPath(file, "attandance", ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if fileUrl != nil {
		data.PhotoUrl = *fileUrl
	}

	allowed := map[string]bool{
		"WFH":    true,
		"VISIT":  true,
		"OFFICE": true,
	}

	if !allowed[data.Activity] {
		return utils.Error(ctx, 400, "activity must be WFH, VISIT, or OFFICE")
	}

	allowedCheckType := map[string]bool{
		"1": true,
		"2": true,
	}

	if !allowedCheckType[data.CheckType] {
		return utils.Error(ctx, 400, "checktype value must be 1 or 2")
	}

	if err := c.repo.Create(&data); err != nil {
		if fileUrl != nil {
			utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Attendance godoc
// @Summary Get all attendance
// @Description Get all attendance with pagination
// @Tags Attendance
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/attendance [get]
func (c *AttendanceController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get My Attendance godoc
// @Summary Get current user attendance
// @Description Get attendance by logged-in user
// @Tags Attendance
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/attendance/me [get]
func (c *AttendanceController) FindByUser(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}
	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data, err := c.repo.FindByUser(userId, &queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}
