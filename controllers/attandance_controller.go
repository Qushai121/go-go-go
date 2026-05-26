package controllers

import (
	"fmt"
	"hrms_go/dto"
	"hrms_go/dto/attandance"
	mappers "hrms_go/mapper"
	"hrms_go/repositories"
	"hrms_go/utils"
	"strconv"
	"strings"

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
// @Param company_code formData string false "Company code"
// @Param office_code formData string true "Office code"
// @Param customer_code formData string false "Customer code"
// @Param logtime formData string true "Log time (YYYY-MM-DD HH:mm:ss or RFC3339)"
// @Param functionno formData int true "Function number"
// @Param action_type formData string false "Action type"
// @Param latitude formData string false "Latitude"
// @Param longitude formData string false "Longitude"
// @Param langtiude formData string false "Longitude alias"
// @Param distance formData string false "Distance"
// @Param imagepath formData file true "Attendance image"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/attendance [post]
func (c *AttendanceController) Create(ctx fiber.Ctx) error {
	request, err := parseAttendanceFormData(ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if request.UserId == "" {
		if userID := ctx.Locals("user_id"); userID != nil {
			request.UserId = fmt.Sprint(userID)
		}
	}

	if request.CreatedBy == "" {
		if employeeNIK := ctx.Locals("employee_nik"); employeeNIK != nil {
			request.CreatedBy = fmt.Sprint(employeeNIK)
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

	employeeNIK := strings.TrimSpace(fmt.Sprint(ctx.Locals("employee_nik")))
	if employeeNIK == "" || employeeNIK == "<nil>" {
		return utils.Error(ctx, 400, "employee_nik is missing from token")
	}
	request.CreatedBy = employeeNIK

	fileUrl, err := utils.SaveFileToCustomPath(file, fmt.Sprintf("employee_absen/%s/foto", employeeNIK), ctx)
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

func parseAttendanceFormData(ctx fiber.Ctx) (attandance.PostAttandanceDto, error) {
	request := attandance.PostAttandanceDto{
		AttendanceId:        ctx.FormValue("attendance_id"),
		UserId:              ctx.FormValue("user_id"),
		CompanyCode:         ctx.FormValue("company_code"),
		OfficeCode:          ctx.FormValue("office_code"),
		CustomerCode:        ctx.FormValue("customer_code"),
		LogTime:             ctx.FormValue("logtime"),
		ActivityType:        ctx.FormValue("activity_type"),
		ActionType:          ctx.FormValue("action_type"),
		Latitude:            ctx.FormValue("latitude"),
		Longitude:           firstFilledFormValue(ctx, "longitude", "langtiude"),
		PresentaseKemiripan: ctx.FormValue("presentase_kemiripan"),
		ImagePath:           ctx.FormValue("imagepath"),
		IsOffline:           ctx.FormValue("is_offline"),
		Distance:            ctx.FormValue("distance"),
		Platforms:           ctx.FormValue("platforms"),
		MaxRadius:           ctx.FormValue("max_radius"),
		ExpandRadius:        ctx.FormValue("expand_radius"),
		ObjectCode:          ctx.FormValue("object_code"),
		CreatedAt:           ctx.FormValue("created_at"),
		UpdatedAt:           ctx.FormValue("updated_at"),
		CreatedBy:           ctx.FormValue("created_by"),
		UpdatedBy:           ctx.FormValue("updated_by"),
	}

	functionNo := ctx.FormValue("functionno")
	if functionNo == "" {
		return request, fmt.Errorf("functionno is required")
	}

	parsedFunctionNo, err := strconv.Atoi(functionNo)
	if err != nil {
		return request, fmt.Errorf("invalid functionno")
	}

	request.FunctionNo = parsedFunctionNo

	return request, nil
}

func firstFilledFormValue(ctx fiber.Ctx, keys ...string) string {
	for _, key := range keys {
		value := ctx.FormValue(key)
		if value != "" {
			return value
		}
	}

	return ""
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
// @Param field_search query string false "Field name for dynamic search"
// @Param activity_type query string false "Filter by activity type"
// @Param logtime query string false "Filter by logtime from date/time"
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
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param activity_type query string false "Filter by activity type"
// @Param logtime query string false "Filter by logtime from date/time"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
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
