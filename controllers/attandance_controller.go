package controllers

import (
	"fmt"
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
// @Param user_id formData string false "User ID (UUID). If empty, taken from auth token."
// @Param company_code formData string true "Company code"
// @Param office_code formData string true "Office code"
// @Param logtime formData string true "Log time (YYYY-MM-DD HH:mm:ss or RFC3339)"
// @Param functionno formData int true "Function number"
// @Param activity_type formData string false "Activity type"
// @Param latitude formData string false "Latitude"
// @Param longitude formData string false "Longitude"
// @Param presentase_kemiripan formData string false "Similarity percentage"
// @Param is_offline formData string false "Offline flag"
// @Param distance formData string false "Distance"
// @Param platforms formData string false "Platform"
// @Param max_radius formData int false "Max radius"
// @Param expand_radius formData int false "Expand radius"
// @Param object_code formData string false "Object code"
// @Param created_by formData string false "Created by"
// @Param updated_by formData string false "Updated by"
// @Param imagepath formData file true "Attendance image"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/attendance [post]
func (c *AttendanceController) Create(ctx fiber.Ctx) error {
	var request attandance.PostAttandanceDto

	if err := ctx.Bind().Body(&request); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if request.UserId == "" {
		if userID := ctx.Locals("user_id"); userID != nil {
			request.UserId = fmt.Sprint(userID)
		}
	}

	if request.CreatedBy == "" {
		if userID := ctx.Locals("user_id"); userID != nil {
			request.CreatedBy = fmt.Sprint(userID)
		}
	}

	if request.ObjectCode == "" {
		request.ObjectCode = "ATTENDANCE"
	}

	file, err := ctx.FormFile("imagepath")
	if err != nil {
		file, err = ctx.FormFile("photo_url")
		if err != nil {
			return utils.Error(ctx, 400, "imagepath file is required")
		}
	}

	fileUrl, err := utils.SaveFileToPath(file, "attandance", ctx)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	request.ImagePath = *fileUrl

	data, err := mappers.ToAttendanceModel(request)
	if err != nil {
		if fileUrl != nil {
			_ = utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 400, err.Error())
	}

	if data.CompanyCode == "" {
		if fileUrl != nil {
			_ = utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 400, "company_code is required")
	}
	if data.OfficeCode == "" {
		if fileUrl != nil {
			_ = utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 400, "office_code is required")
	}
	if data.FunctionNo <= 0 {
		if fileUrl != nil {
			_ = utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 400, "functionno must be greater than 0")
	}

	if err := c.repo.Create(&data); err != nil {
		if fileUrl != nil {
			_ = utils.RemoveFileFromPath(*fileUrl)
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

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
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
	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data, err := c.repo.FindByUser(userId, &queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}
